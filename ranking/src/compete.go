package compete

import (
    "ranking/config"
    "ranking/models"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
    "context"
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

    options := types.ImageBuildOptions{
        Tags: []string{config.GameTag},
        SuppressOutput: true,                           
        Dockerfile: "Dockerfile",           
    }
    resp, err := Dock_cli.ImageBuild(context.Background(), file, options)
    defer resp.Body.Close()

    return err
}


