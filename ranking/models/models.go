package models

import (
    "database/sql"
    "ranking/config"
    "github.com/mafredri/go-trueskill"
    "sync"
    "strings"
    "fmt"
)

var DB *sql.DB
var mutex sync.Mutex

type Agent struct {
    Id int64 
    UserId int 
    FileName string 
    Image string `json:"stream"`
    Raiting float64 
    Sigma float64
    Broken bool
}

//type PlayedMatch struct {
//    Winner *Agent
//    Looser *Agent
//    Draw bool
//}

func CreateAgentDB(data *Agent, owner Player) error {
    data.Raiting = config.DefMu
    data.Sigma = config.DefSig
    err := DB.QueryRow("INSERT INTO submissions (user_id, file_name,  container_id, raiting, sigma) VALUES ($1, $2, $3, $4, $5) RETURNING id", owner.Id, data.FileName, data.Image, data.Raiting, data.Sigma).Scan(&data.Id);
    data.UserId = owner.Id
    return err;
}


func GetClosest(current *Agent, skip []int64) (Agent, bool) {
    com := strings.Trim(strings.Replace(fmt.Sprint(skip), " ", ",", -1), "[]")
    rows, _:= DB.Query("SELECT * FROM submissions where user_id!=? and broken=0 and id not in (?)", current.UserId, com)
    defer rows.Close()
    var res Agent
    best_quality := -1.0 
    curr_skill := trueskill.NewPlayer(current.Raiting, current.Sigma)
    for rows.Next() {
        var agent Agent
        var broken int
        rows.Scan(&agent.Id, &agent.UserId, &agent.Image, &agent.Raiting, &agent.Sigma, &broken)
        cand_skill := trueskill.NewPlayer(agent.Raiting, agent.Sigma)
        quality := config.TsConfig().MatchQuality([]trueskill.Player{curr_skill, cand_skill})
        if quality > best_quality {
            res = agent
            best_quality = quality
        }
    }
    return res, best_quality > -1

}

func SetBroken(agent *Agent) {
   DB.Exec("update submissions set broken=1 where id=?", agent.Id)
}

func reloadAgent(agent *Agent) error {
    return DB.QueryRow("select raiting, sigma from submissions where id=?", agent.Id).Scan(&agent.Raiting, &agent.Sigma)
}

func updateAgent(ndata *Agent, pl *trueskill.Player) error {
    _, err := DB.Exec("update submissions set raiting=?, sigma=? where id=?", pl.Mu(), pl.Sigma(), ndata.Id)
    return err
}

func RecordResult(winner *Agent, looser *Agent, draw bool) {
    mutex.Lock()
    defer mutex.Unlock()
    reloadAgent(winner)
    reloadAgent(looser)
    pl1 := trueskill.NewPlayer(winner.Raiting, winner.Sigma)
    pl2 := trueskill.NewPlayer(looser.Raiting, looser.Sigma)
    newSkills, _ := config.TsConfig().AdjustSkills([]trueskill.Player{pl1, pl2}, draw)
    updateAgent(winner, &newSkills[0])
    updateAgent(looser, &newSkills[1])
}

func GetAgentsN() (int, error) {
    var count int
    err := DB.QueryRow("select count(id) from submissions").Scan(&count)
    return count, err
}

type Player struct {
    Id int
    Name string
    Raiting float64
    Agents int
}

func GetOrCreatePlayer(user_name string) (Player, error) {
    // TODO: defence from sql injections.
    var user Player
    err := DB.QueryRow("insert into players (user_name) values ('"+user_name+"')  on conflict do nothing; select * from players where user_name='"+user_name+"'").Scan(&user.Id, &user.Name);
    return user, err
}

func GetLeaderboard() []Player {
    rows, _:= DB.Query("SELECT * FROM players")
    defer rows.Close()
    var res []Player
    for rows.Next() {
        var player Player
        rows.Scan(&player.Id, &player.Name);
        player.Raiting = 173
        player.Agents = 7
        res = append(res, player);
    }
    return res
}
