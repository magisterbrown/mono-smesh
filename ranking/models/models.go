package models

import (
    "database/sql"
    "ranking/config"
    "math"
    "github.com/mafredri/go-trueskill"
    "sync"
    "strings"
    "fmt"
)

var DB *sql.DB
var mutex sync.Mutex

type Agent struct {
    Id int64
    UserId int //TODO: fill it
    Image string `json:"stream"`
    Raiting float64
    Sigma float64
}

//type PlayedMatch struct {
//    Winner *Agent
//    Looser *Agent
//    Draw bool
//}

func SaveAgent(data *Agent, owner Player) error {
    data.Raiting = config.DefMu
    data.Sigma = config.DefSig
    res, err := DB.Exec("INSERT INTO submissions (user_id, container_id, raiting, sigma) VALUES (?, ?, ?, ?)", owner.Id, data.Image, data.Raiting, data.Sigma);
    if err != nil {
        return err
    }
    idx, err := res.LastInsertId();
    data.Id = idx
    data.UserId = owner.Id
    return err;
}


func GetClosest(current *Agent, skip []int64) (Agent, bool) {
    com := strings.Trim(strings.Replace(fmt.Sprint(skip), " ", ",", -1), "[]")
    rows, _:= DB.Query("SELECT * FROM submissions where user_id!=? and broken=0 and id not in (?)", current.UserId, com)
    defer rows.Close()
    var res Agent
    diff := math.Inf(1)
    hasRows := false
    for rows.Next() {
        var agent Agent
        var broken int
        rows.Scan(&agent.Id, &agent.UserId, &agent.Image, &agent.Raiting, &agent.Sigma, &broken)
        ndiff := math.Abs(current.Raiting - agent.Raiting)
        if ndiff<diff {
            res = agent
            ndiff = diff
            hasRows = true
        }
    }
    return res, hasRows

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
}

func GetPlayer(token string) (Player, error) {
    // TODO: defence from sql injections.
    var user Player
    err := DB.QueryRow("select * from players where id in (select user_id from sessions where token='"+token+"')").Scan(&user.Id, &user.Name);
    return user, err
}

func GetLeaderboard() []Player {
    rows, _:= DB.Query("SELECT * FROM players")
    defer rows.Close()
    var res []Player
    for rows.Next() {
        var player Player
        rows.Scan(&player.Id, &player.Name);
        res = append(res, player);
    }
    return res
}
