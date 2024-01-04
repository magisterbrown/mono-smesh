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
        winner, looser, draw := MatchShuffle(&agent, &competitor)
        models.RecordResult(winner, looser, draw)
        extra_competitor, hasRows := models.GetClosest(&competitor, []int64{agent.Id})
        if hasRows {
            go func() {
                winner, looser, draw := MatchShuffle(&extra_competitor, &competitor)
                models.RecordResult(winner, looser, draw)
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


