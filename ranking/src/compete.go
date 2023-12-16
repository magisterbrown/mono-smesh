package compete

import (
    "ranking/config"
    "ranking/models"
    "fmt"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/strslice"
    "context"
    "encoding/json"
)

func ScheduleNGames(agent models.Agent, nGames int) {
    played := []int64{}
    for i:=0; i<nGames;i++ {
        competitor, hasRows := models.GetClosest(&agent, played)
        if !hasRows {
            return 
        }
        played = append(played, competitor.Id)
        winner, err := Match(&agent, &competitor)
        if err != nil{
            models.SetBroken(winner)
            if winner.Id == agent.Id {
                return 
            }
        } else {
            //TODO: support draws
            if winner == nil {
                //models.RecordResult(&agent, winner, true)
            }else if(winner.Id == agent.Id){
                //models.RecordResult(&agent, winner, false)
            } else {
                //models.RecordResult(winner, &agent, false)
            }
        }
        //comparable := make([]trueskill.Player, len(agents))
        //for i, ag := range agents{
        //    comparable[i] = trueskill.NewPlayer(int(ag.Id), ag.Raiting, ag.Sigma) 
        //}
    }
}

func InitGame() error {
    file, err := archive.TarWithOptions(config.GameFolder, &archive.TarOptions{IncludeFiles: []string{"Dockerfile", "play.py", "requirements.txt"}})
    if err != nil {
	    return err
	}
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
	    return err
	}

    options := types.ImageBuildOptions{
        Tags: []string{config.GameTag},
        SuppressOutput: true,                           
        Dockerfile: "Dockerfile",           
    }
    resp, err := cli.ImageBuild(context.Background(), file, options)
    defer resp.Body.Close()

    return err
}

func startContainer(cli *client.Client, image string, command string) (types.HijackedResponse, string, error) {

    conf := container.Config{Image: image, AttachStdin: true, AttachStdout: true, AttachStderr: true, Tty: true, OpenStdin: true, StdinOnce: true, Cmd:  strslice.StrSlice{command}}
    hijack := types.HijackedResponse{}
    playercontID := ""
    err := *new(error)

    for range []int{1} {
        playercont, err := cli.ContainerCreate(context.Background(), &conf, nil, nil, nil, "")
        if err != nil {
            break
        }
        playercontID = playercont.ID
        hijack, err = cli.ContainerAttach(context.Background(), playercont.ID, types.ContainerAttachOptions{Stream:true, Stdout:true, Stdin:true, Stderr:true})
        if err != nil {
            break
        }
        err = cli.ContainerStart(context.Background(), playercont.ID, types.ContainerStartOptions{})
        if err != nil {
            break
        }
    }

    return hijack, playercontID, err
}


func Match(player1 *models.Agent, player2 *models.Agent) (*models.Agent, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        panic(err)
    }
    
    //Start game container
    hijack, gamecontID, err := startContainer(cli, config.GameTag, "")
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
                output, containerID, err := startContainer(cli, agents[agent_idx].Image, string(buf))
                if err != nil{
                    break gameLoop
                }
                line, isPrefix, err := output.Reader.ReadLine()
                if err != nil || isPrefix {
                    err = fmt.Errorf("Failed to get full line from user %w", err)
                }
                output.Close()
                err = cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
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
    errclean := cli.ContainerRemove(context.Background(), gamecontID, types.ContainerRemoveOptions{Force: true})
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
