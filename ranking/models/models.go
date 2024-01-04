package models

import (
    "database/sql"
    "ranking/config"
    "github.com/lib/pq"
    "github.com/mafredri/go-trueskill"
    "time"
)

var DB *sql.DB

type Agent struct {
    Id int64 
    UserId int 
    FileName string 
    Image string `json:"stream"`
    Raiting float64 
    Sigma float64
    Broken bool
    CreatedAt time.Time
    
}

func CreateAgentDB(data *Agent, owner Player) error {
    data.Raiting = config.DefMu
    data.Sigma = config.DefSig
    err := DB.QueryRow("INSERT INTO submissions (user_id, file_name,  container_id, raiting, sigma) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at", owner.Id, data.FileName, data.Image, data.Raiting, data.Sigma).Scan(&data.Id, &data.CreatedAt);
    data.UserId = owner.Id
    return err;
}


func GetClosest(current *Agent, skip []int64) (Agent, bool) {
    //skip = append(skip, -1)
    //com := strings.Trim(strings.Replace(fmt.Sprint(skip), " ", ",", -1), "[]")
    rows, err:= DB.Query("SELECT * FROM submissions where user_id!=$1 and broken=0 and id!=ALL($2)", current.UserId, pq.Array(skip))
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    var res Agent
    best_quality := -1.0 
    curr_skill := trueskill.NewPlayer(current.Raiting, current.Sigma)

    for rows.Next() {

        var agent Agent
        var broken int
        rows.Scan(&agent.Id, &agent.UserId, &agent.FileName, &agent.Image, &agent.Raiting, &agent.Sigma, &broken, &agent.CreatedAt)
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
    _, err := DB.Exec("update submissions set broken=1 where id=$1", agent.Id)
   if err != nil {
        panic(err)
   }
}

func ReloadAgent(agent *Agent) {
    err := DB.QueryRow("select raiting, sigma from submissions where id=$1", agent.Id).Scan(&agent.Raiting, &agent.Sigma)
    if err != nil {
        panic(err)
    }
}

func UpdateAgent(ndata *Agent, pl *trueskill.Player) {
    _, err := DB.Exec("update submissions set raiting=$1, sigma=$2 where id=$3", pl.Mu(), pl.Sigma(), ndata.Id)
    if err != nil {
        panic(err)
    }
}

func StoreMatch(history string) int {
    var id int;
    err := DB.QueryRow("INSERT INTO matches (recording) VALUES ($1) RETURNING id", history).Scan(&id)
    if err != nil {
        panic(err)
    }
    return id
}

func StoreSeating(matchId int, agent *Agent, pl *trueskill.Player, spot string) {
    _, err := DB.Exec("INSERT INTO seating (match_id, submission_id, change, spot) VALUES ($1, $2, $3, $4)", matchId, agent.Id,  pl.Mu() - agent.Raiting, spot)
    if err != nil {
        panic(err)
    }
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
