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
    Args map[string]interface{}
}


func Match(player1 *models.Agent, player2 *models.Agent) (*models.Agent, error) {
    sock_name := "/" + uuid.New().String() + ".sock";
    socket, err := net.Listen("unix", config.SockerVolumePath + sock_name)
    if err != nil {
        panic(err)
    }
    defer socket.Close()
    fmt.Println(sock_name);

    conf := container.Config{Image: config.GameTag, Cmd: strslice.StrSlice{"/sockets" + sock_name}}
    hostConf := container.HostConfig{Binds: []string{config.SocketVolumeName + ":/sockets"}}
    _ = hostConf
    _ = conf
    gamecont, err := Dock_cli.ContainerCreate(context.Background(), &conf, &hostConf, nil, nil, "")
    if err != nil {
        panic(err)
    }
    err = Dock_cli.ContainerStart(context.Background(), gamecont.ID, types.ContainerStartOptions{})
    if err != nil {
        panic(err)
    }
    
    //Event loop of the game
    for {
        conn, err := socket.Accept()
        if err != nil {
            panic(err)
        }
        buf := make([]byte, 4096)
        n, err := conn.Read(buf)
        req := Request{}
        if err = json.Unmarshal(buf[:n], &req); err != nil {
            panic(err)
        }
        var command string 
        for key, value := range req.Args {
            command += fmt.Sprintf(" --%s %v", key, value)
        }
        fmt.Println(command)
        conn.Close()
        break
    }
    return player1, nil
    
    //Start game container
    hijack, gamecontID, err := startContainer(config.GameTag, "")
    if err != nil {
        panic(err)
    }
    defer hijack.Close()


    agents := map[string]*models.Agent{
            "player_0": player1,
            "player_1": player2,
    }

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
        return agents[agent_idx], err
    }
    if winner == "" {
        return nil, nil
    }
    
    return agents[winner], nil
    }
