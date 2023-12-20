package main

import (
    "fmt"
    _ "math"
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
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD": 
            res := models.GetLeaderboard()
            json.NewEncoder(w).Encode(res)
        case "POST":
            
            //Get from content

            // TODO: add form size limit to config
            req.ParseMultipartForm(11 << 20) 
            player, err := models.GetOrCreatePlayer(req.Header.Get("User-Name"));
            _ = player;
            if err != nil {
            	http.Error(w, "Idk how in the world you got this error. Try resubmitting or new acount", http.StatusBadRequest)
            	return
            }
            file, header, err := req.FormFile("file")
            if err != nil {
            	http.Error(w, "Could not read file named file from a from", http.StatusBadRequest)
            	return
            }
            defer file.Close()
                                   
            // Create docker image
            options := types.ImageBuildOptions{
                Tags: []string{petname.Generate(2, "-")+":player"},
                SuppressOutput: true,                           
                Dockerfile: "submission/Dockerfile",           
            }                                                 

            // TODO: limit build time
            resp, err := compete.Dock_cli.ImageBuild(context.Background(), file, options)
            if err != nil {
                http.Error(w, "Failed to build and image: "+err.Error(), http.StatusBadRequest)
		        return
	        }
            defer resp.Body.Close()

            //Crete Agent
            var submission models.Agent
            submission.FileName = header.Filename
            submission.Broken = false;

            json.NewDecoder(resp.Body).Decode(&submission)
            submission.Image = strings.TrimSuffix(strings.TrimPrefix(submission.Image, "sha256:"), "\n")

            // TODO: play one match against itself
            _, err = compete.Match(&submission, &submission);
            //if(err != nil){
            //	http.Error(w, "Agent does not play by the rules", http.StatusBadRequest)
            //	return
            //}
            err = models.CreateAgentDB(&submission, player)
            if err != nil{
                http.Error(w, "Agent works correct but did not save because: "+err.Error(), http.StatusBadRequest)
		        return
            }

            // TODO: Schedule games
            //competitors, err := models.GetAgentsN()
            //if err != nil{
            //    panic(err)
            //}
            //matches := int(math.Ceil(2*math.Log2(float64(competitors))))
            //go compete.ScheduleNGames(submission, matches)
           
            w.WriteHeader(200);
            json.NewEncoder(w).Encode(submission) //TODO: handle serializaiton errors
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

	compete.Dock_cli, err = client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        panic(err)
    }

    err = compete.InitGame()
    if err != nil {
        panic(err)
    }

    fmt.Println("Initialized all systems")
    http.HandleFunc("/api/leaderboard", getLeaderboard)
    http.ListenAndServe(":5000", nil)
}
