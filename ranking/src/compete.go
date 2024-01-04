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
        outcome := MatchShuffle(&agent, &competitor)
        RecordResult(&outcome)
        extra_competitor, hasRows := models.GetClosest(&competitor, []int64{agent.Id})
        if hasRows {
            go func() {
                outcome := MatchShuffle(&extra_competitor, &competitor)
                RecordResult(&outcome)
            }()
        }
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


