package compete

import (
    "fmt"
    "net"
    "context"
    "time"
    "encoding/json"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/strslice"
    "github.com/google/uuid"
    "ranking/config"
    "ranking/models"
    "math/rand"
)

var Dock_cli *client.Client;

type Request struct {
    Type string  
    Agent string

    //Optional fields
    Args map[string]interface{}
    Broken string
}

func MatchShuffle(player1 *models.Agent, player2 *models.Agent) (*models.Agent, *models.Agent) {
    rand.Seed(time.Now().UnixNano())
    if(rand.Intn(2) == 1){
        return Match(player1, player2)
    }
    return Match(player2, player1)
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
        n, err := conn.Read(buf)
        req = Request{}
        if err = json.Unmarshal(buf[:n], &req); err != nil {
            panic(err)
        }
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
}
