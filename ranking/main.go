package main

import (
    "fmt"
    "math"
    "net/http"
    "io/ioutil"
    "database/sql"
    "encoding/json"
    "ranking/models"
    "ranking/src"
    "ranking/config"
    _ "github.com/lib/pq" 
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/dustinkirkland/golang-petname"
    "context"
    "strings"
)

func getLeaderboard(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch req.Method {
        case "OPTIONS":
            w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
            //w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD": 
            res := models.GetLeaderboard()
            json.NewEncoder(w).Encode(res)
        case "POST":
            req.ParseMultipartForm(11 << 20) 

            //Authentication
            player, err := models.GetPlayer(req.Header.Get("Authorization"));
            if err != nil {
            	http.Error(w, "Authorization failed", http.StatusUnauthorized)
                return 
            }

            //Processing form
            file, _, err:= req.FormFile("submission")
            if err != nil {
            	http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
            	return
            }
            defer file.Close()
            
            // Create docker image
	        cli, err := client.NewClientWithOpts(client.FromEnv)
            options := types.ImageBuildOptions{
                Tags: []string{petname.Generate(2, "-")+":player"},
                SuppressOutput: true,                           
                Dockerfile: "submission/Dockerfile",           
            }                                                 
            resp, err := cli.ImageBuild(context.Background(), file, options)
            
            if err != nil {
		        fmt.Println("Error:", err)
		        return
	        }
            defer resp.Body.Close()

            var submission models.Agent
            json.NewDecoder(resp.Body).Decode(&submission)
            submission.Image = strings.TrimSuffix(strings.TrimPrefix(submission.Image, "sha256:"), "\n")
            _, err = compete.Match(&submission, &submission);
            if(err != nil){
            	http.Error(w, "Agent does not play by the rules", http.StatusBadRequest)
            	return
            }
            err = models.SaveAgent(&submission, player)
            if err != nil{
                panic(err)
            }

            // TODO: Schedule games
            competitors, err := models.GetAgentsN()
            if err != nil{
                panic(err)
            }
            matches := int(math.Ceil(2*math.Log2(float64(competitors))))
            go compete.ScheduleNGames(submission, matches)

            json.NewEncoder(w).Encode(map[string]string{"status":"ok", "raiting": "600"})
        default:
            w.WriteHeader(404)
    }
}

func main() {
    var err error
    models.DB, err = sql.Open("postgres", config.DBurl) 
    if err != nil {
        panic(err)
    }

    file, err := ioutil.ReadFile("./schema.sql")
    if err != nil {
        panic(err)
    }
    _, err = models.DB.Exec(string(file))
    if err != nil {
        panic(err)
    }

    //ct, _ := models.GetAgentsN()
    fmt.Println("Started")


    err = compete.InitGame()

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/api/leaderboard", getLeaderboard)
    http.ListenAndServe(":5000", nil)
}
