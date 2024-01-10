package compete

import (
    "fmt"
    "net"
    "context"
    "time"
    "sync"
    "encoding/json"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/strslice"
    "github.com/google/uuid"
    "github.com/mafredri/go-trueskill"
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
    History string
}

type Outcome struct {
    Winner *models.Agent
    Looser *models.Agent
    Draw bool
    Recording string
    Seating map[*models.Agent]string
}

func MatchShuffle(player1 *models.Agent, player2 *models.Agent) Outcome {
    rand.Seed(time.Now().UnixNano())
    if(rand.Intn(2) == 1){
        return Match(player1, player2)
    }
    return Match(player2, player1)
}


func Match(player1 *models.Agent, player2 *models.Agent) Outcome {
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
            "": nil,
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
    
    out := new(Outcome)
    out.Draw = false
    out.Recording = req.History
    out.Seating = make(map[*models.Agent]string)

    for k, v := range agents{
        out.Seating[v] = k 
    }

    if req.Agent == "" {
        out.Draw = true
    }

    if agents[req.Agent] == player1 {
        out.Winner = player1
        out.Looser = player2
    } else {
        out.Winner = player2
        out.Looser = player1
    }

    return *out
}

var mutex sync.Mutex

func RecordResult(out *Outcome) {
    mutex.Lock()
    defer mutex.Unlock()
    models.ReloadAgent(out.Winner)
    models.ReloadAgent(out.Looser)
    pl1 := trueskill.NewPlayer(out.Winner.Raiting, out.Winner.Sigma)
    pl2 := trueskill.NewPlayer(out.Looser.Raiting, out.Looser.Sigma)
    newSkills, _ := config.TsConfig().AdjustSkills([]trueskill.Player{pl1, pl2}, out.Draw)
    matchId := models.StoreMatch(out.Recording)
    if(out.Draw) {
        models.StoreSeating(matchId, out.Winner, &newSkills[0], out.Seating[out.Winner], "Draw"); 
        models.StoreSeating(matchId, out.Looser, &newSkills[1], out.Seating[out.Looser], "Draw"); 
    } else{
        models.StoreSeating(matchId, out.Winner, &newSkills[0], out.Seating[out.Winner], "Win"); 
        models.StoreSeating(matchId, out.Looser, &newSkills[1], out.Seating[out.Looser], "Loose"); 
    }
    models.UpdateAgent(out.Winner, &newSkills[0])
    models.UpdateAgent(out.Looser, &newSkills[1])

    //_, err := DB.Exec("INSERT INTO matches (subm_1, subm_2, subm_1_change, subm_2_change, recording) VALUES ($1, $2, $3, $4, $5)", winner.Id, looser.Id, newSkills[0].Mu() - winner.Raiting, newSkills[1].Mu() - winner.Raiting, "{}")
    //if err != nil {
    //    panic(err)
    //}

}
