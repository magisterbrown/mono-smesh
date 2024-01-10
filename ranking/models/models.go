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

type SpotPlayer struct {
    Id int64
    Change float64
    Spot string
    Status string
    UserName string
}

type Submission struct {
    Id int64
    Seating []SpotPlayer
}

func CreateAgentDB(data *Agent, owner Player) error {
    data.Raiting = config.DefMu
    data.Sigma = config.DefSig
    err := DB.QueryRow("INSERT INTO submissions (user_id, file_name,  container_id, raiting, sigma) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at", owner.Id, data.FileName, data.Image, data.Raiting, data.Sigma).Scan(&data.Id, &data.CreatedAt);
    data.UserId = owner.Id
    return err;
}

func GetSubmissions(id int64) ([]Submission, error) {
    rows, err := DB.Query("select m.id from matches m join seating s on m.id=s.match_id join submissions su on s.submission_id=su.id where su.id=$1", id)
    if err != nil {
        return nil, err
    }
    var subms []Submission
    for rows.Next() {
        var subm Submission
        subm.Seating = []SpotPlayer{}
        rows.Scan(&subm.Id)
        spotRows, err := DB.Query("select se.id, se.change, se.spot, se.status, p.user_name from seating se join submissions su on se.submission_id=su.id join players p on su.user_id=p.id where match_id=$1", subm.Id)
        if err != nil {
            return subms, err 
        }
        for spotRows.Next() {
            var spot SpotPlayer
            spotRows.Scan(&spot.Id, &spot.Change, &spot.Spot, &spot.Status, &spot.UserName)
            subm.Seating = append(subm.Seating, spot)
        }
        subms = append(subms, subm);
    }
    return subms, nil
}

func GetRecording(id int64) (string, error) {
    var recording string
    err := DB.QueryRow("select recording from matches where id=$1", id).Scan(&recording);
    return recording, err
}

func GetUserAgents(user_name string) ([]Agent, error) {
    //TODO: protect from injections
    rows, err := DB.Query("SELECT s.id, s.user_id, s.file_name, s.container_id, s.raiting, s.sigma, s.broken, s.created_at FROM submissions s JOIN players p ON s.user_id=p.id WHERE p.user_name=$1 order by s.created_at DESC", user_name)
    if err != nil {
       return nil, err
    }
    var agents []Agent
    for rows.Next() {
        var agn Agent
        rows.Scan(&agn.Id, &agn.UserId, &agn.FileName, &agn.Image, &agn.Raiting, &agn.Sigma, &agn.Broken, &agn.CreatedAt)
        agents = append(agents, agn)
    }
    if err = rows.Err(); err != nil {
        return agents, err
    }
    return agents, nil
}

func GetClosest(current *Agent, skip []int64) (Agent, bool) {
    rows, err:= DB.Query("SELECT * FROM submissions where user_id!=$1 and broken='f' and id!=ALL($2)", current.UserId, pq.Array(skip))
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
    _, err := DB.Exec("update submissions set broken='t' where id=$1", agent.Id)
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

func StoreSeating(matchId int, agent *Agent, pl *trueskill.Player, spot string, status string) {
    _, err := DB.Exec("INSERT INTO seating (match_id, submission_id, change, spot, status) VALUES ($1, $2, $3, $4, $5)", matchId, agent.Id,  pl.Mu() - agent.Raiting, spot, status)
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
    BestAgentId int
}

func GetOrCreatePlayer(user_name string) (Player, error) {
    // TODO: defence from sql injections.
    var user Player
    err := DB.QueryRow("insert into players (user_name) values ('"+user_name+"')  on conflict do nothing; select * from players where user_name='"+user_name+"'").Scan(&user.Id, &user.Name);
    return user, err
}

func GetPlayerId(user_name string) (int, error) {
    var id int;
    err := DB.QueryRow("SELECT id FROM players WHERE user_name=$1", user_name).Scan(&id);
    return id, err
}

func GetLeaderboard() []Player {
    rows, _:= DB.Query("with counted as (select p.id, p.user_name, MAX(s.raiting), COUNT(s.id) from players p join submissions s on p.id=s.user_id group by p.id) select c.*, s.id max_sub_id from counted c join submissions s on s.user_id=c.id and s.raiting=c.max order by max DESC")
    defer rows.Close()
    var res []Player
    for rows.Next() {
        var player Player
        rows.Scan(&player.Id, &player.Name, &player.Raiting, &player.Agents, &player.BestAgentId);
        res = append(res, player);
    }
    return res
}
