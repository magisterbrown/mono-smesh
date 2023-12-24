package compete

import (
    "fmt"
    "net"
    "context"
    "encoding/json"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/strslice"
    "github.com/google/uuid"
    "ranking/config"
    "ranking/models"
)

var Dock_cli *client.Client;

func startContainer(image string, command string) (types.HijackedResponse, string, error) {

    conf := container.Config{Image: image, AttachStdin: true, AttachStdout: true, AttachStderr: true, Tty: true, OpenStdin: true, StdinOnce: true, Cmd:  strslice.StrSlice{command}}
    hijack := types.HijackedResponse{}
    playercontID := ""
    err := *new(error)

    for range []int{1} {
        playercont, err := Dock_cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")
        if err != nil {
            break
        }
        playercontID = playercont.ID
        hijack, err = Dock_cli.ContainerAttach(context.Background(), playercont.ID, types.ContainerAttachOptions{Stream:true, Stdout:true, Stdin:true, Stderr:true})
        if err != nil {
            break
        }
        err = Dock_cli.ContainerStart(context.Background(), playercont.ID, types.ContainerStartOptions{})
        if err != nil {
            break
        }
    }

    return hijack, playercontID, err
}

type Request struct {
    Type string  
    Agent string

    //Optional fields
    Args map[string]interface{}
    Broken string
}


func Match(player1 *models.Agent, player2 *models.Agent) (*models.Agent, *models.Agent) {
    sock_name := "/" + uuid.New().String() + ".sock";
    socket, err := net.Listen("unix", config.SockerVolumePath + sock_name)
    if err != nil {
        panic(err)
    }
    defer socket.Close()

    conf := container.Config{Image: config.GameTag, Cmd: strslice.StrSlice{"/sockets" + sock_name}}
    hostConf := container.HostConfig{Binds: []string{config.SocketVolumeName + ":/sockets"}}
    gamecont, err := Dock_cli.ContainerCreate(context.Background(), &conf, &hostConf, nil, nil, "")
    if err != nil {
        panic(err)
    }
    err = Dock_cli.ContainerStart(context.Background(), gamecont.ID, types.ContainerStartOptions{})
    if err != nil {
        panic(err)
    }
    
    agents := map[string]*models.Agent{
            "player_0": player1,
            "player_1": player2,
    }

    req := Request{}
    //Event loop of the game
    for {
        conn, err := socket.Accept()
        if err != nil {
            panic(err)
        }
        buf := make([]byte, 4096)
        fmt.Println("waiting")
        n, err := conn.Read(buf)
        fmt.Println("income")
        req = Request{}
        if err = json.Unmarshal(buf[:n], &req); err != nil {
            panic(err)
        }
        fmt.Println(req.Agent)
        if req.Type == "done" {
            break
        }
        go func(req Request, conn net.Conn) {
            defer conn.Close()
            var command string 
            for key, value := range req.Args {
                command += fmt.Sprintf(" --%s %v", key, value)
            }
            agent := agents[req.Agent]

            //Running agent
            err = func() error {
                conf := container.Config{Image: agent.Image, Tty: true,  Cmd:  strslice.StrSlice{command}}
                cont, err := Dock_cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")
                if(err != nil) {
                    return err
                }

                hijack, err := Dock_cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream:true, Stdout:true})

                if(err != nil) {
                    return err
                }

                if err = Dock_cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
                    return err
                }
                buf, _, err:= hijack.Reader.ReadLine()
                if(err != nil) {
                    return err
                }
                fmt.Println("stepp")
                fmt.Println(string(buf))

                if err = Dock_cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
                    return err
                }
                conn.Write(buf)
                return nil
            }()
            if err != nil{
                conn.Write([]byte("broken"))
            }

        }(req, conn)
    }
    err = Dock_cli.ContainerRemove(context.Background(), gamecont.ID, types.ContainerRemoveOptions{Force: true})
    if(err != nil) {
        panic(err)
    }


    if req.Broken != "" {
       agents[req.Broken].Broken = true
    }

    if req.Agent == "" {
        return nil, nil
    }

    if agents[req.Agent] == player1 {
        return player1, player2
    }

    return player2, player1

    //Start game container
    hijack, gamecontID, err := startContainer(config.GameTag, "")
    if err != nil {
        panic(err)
    }
    defer hijack.Close()


    

    //Play game
    var message map[string]interface{}
    agent_idx := ""
    winner := ""
    gameLoop:
    for {
        buf, isPrefix, err := hijack.Reader.ReadLine()
        if err != nil || isPrefix {
            err = fmt.Errorf("Failed to get full line from user %w", err)
        }
        err = json.Unmarshal(buf,&message)
        if err != nil{
            break gameLoop
        }
        intype, ok := message["type"].(string)
        if !ok {
            err = fmt.Errorf("Not type from game")
            break gameLoop
        }
        switch intype {
            case "move":{
                agent_idx, ok = message["agent"].(string)
                if !ok {
                    err = fmt.Errorf("No agent id")
                    break gameLoop
                }
                output, containerID, err := startContainer(agents[agent_idx].Image, string(buf))
                if err != nil{
                    break gameLoop
                }
                line, isPrefix, err := output.Reader.ReadLine()
                if err != nil || isPrefix {
                    err = fmt.Errorf("Failed to get full line from user %w", err)
                }
                output.Close()
                err = Dock_cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
                if err != nil{
                    break gameLoop
                }
                hijack.Conn.Write(append(line,[]byte("\n")...))
            }
            case "result":{
                winner, _ = message["winner"].(string)
                break gameLoop
            }
            default:
                err = fmt.Errorf("Unexpected json")
                break gameLoop
        }
    }

    //Cleanup game
    errclean := Dock_cli.ContainerRemove(context.Background(), gamecontID, types.ContainerRemoveOptions{Force: true})
    if errclean != nil {
        panic(errclean)
    }
    if agent_idx == "" {
        panic(fmt.Errorf("Game did not set the player"))
    }
    if err != nil {
        return agents[agent_idx], nil
    }
    if winner == "" {
        return nil, nil
    }
    
    return agents[winner], nil
    }
