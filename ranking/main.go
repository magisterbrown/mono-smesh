package main

import (
    "fmt"
    "regexp"
    "io"
    "math"
    "strconv"
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
)
func getMyName(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    switch req.Method {
        case "OPTIONS":
            w.Header().Set("Allow", "OPTIONS, GET, HEAD")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD":
            w.Write([]byte(req.Header.Get("User-Name")))
        default:
            w.WriteHeader(404);
    }
}
func getRecording(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    switch req.Method {
        case "OPTIONS":
            w.Header().Set("Allow", "OPTIONS, GET, HEAD")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD":
            idx, ok := req.URL.Query()["id"]
            if !ok || len(idx) != 1 {
            	http.Error(w, "Sepcify one ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            idntx, err := strconv.Atoi(idx[0]);
            if err != nil {
            	http.Error(w, "Sepcify one ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            recording, err := models.GetRecording(int64(idntx))
            if err != nil {
            	http.Error(w, "Sepcify valid ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            w.Write([]byte(recording))
        default:
            w.WriteHeader(404);
    }
}

func getSubmissions(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch req.Method {
        case "OPTIONS":
            w.Header().Set("Allow", "OPTIONS, GET, HEAD")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD":
            userName, ok := req.URL.Query()["user_name"]
            if !ok || len(userName) != 1 {
            	http.Error(w, "Sepcify one ?user_name=<name>", http.StatusBadRequest)
            	return
            }
            agents, err := models.GetUserAgents(userName[0])
            if err != nil {
                http.Error(w, "Try better user_name", http.StatusBadRequest)
                return 
            }
            json.NewEncoder(w).Encode(agents)
        default:
            w.WriteHeader(404);
    }

}

func getMatches(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch req.Method {
        case "OPTIONS":
            w.Header().Set("Allow", "OPTIONS, GET, HEAD")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            w.WriteHeader(http.StatusOK)
        case "GET", "HEAD":
            idx, ok := req.URL.Query()["id"]
            if !ok || len(idx) != 1 {
            	http.Error(w, "Sepcify one ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            idntx, err := strconv.Atoi(idx[0]);
            if err != nil {
            	http.Error(w, "Sepcify one ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            submissions, err := models.GetSubmissions(int64(idntx))
            if err != nil {
            	http.Error(w, "Sepcify valid ?id=<submission_id>", http.StatusBadRequest)
                return 
            }
            json.NewEncoder(w).Encode(submissions)
        default:
            w.WriteHeader(404);
    }
}

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
                Dockerfile: "Dockerfile",           
            }                                                 

            // TODO: limit build time
            resp, err := compete.Dock_cli.ImageBuild(context.Background(), file, options)
            if err != nil {
                http.Error(w, "Falied before building: "+err.Error(), http.StatusBadRequest)
		        return
	        }
            defer resp.Body.Close()

            buf := make([]byte, 4096)
            n, err := resp.Body.Read(buf)
            if err != nil && err != io.EOF{
                http.Error(w, "Error when reading build result: "+err.Error(), http.StatusBadRequest)
		        return
	        }
            re := regexp.MustCompile(`^{"stream":"sha256:(\w*)\\n"}`)
            idxMatch := re.FindSubmatch(buf[:n])
            if len(idxMatch) != 2{
                http.Error(w, "Image build failed: "+string(buf[:n]), http.StatusBadRequest)
		        return
            }

            //Crete Agent
            var submission models.Agent
            submission.FileName = header.Filename
            submission.Broken = false
            submission.Image = string(idxMatch[1])

            compete.Match(&submission, &submission);
            if(submission.Broken){
            	http.Error(w, "Agent does not play by the rules", http.StatusBadRequest)
            	return
            }
            
            // Save agent
            err = models.CreateAgentDB(&submission, player)
            if err != nil{
                http.Error(w, "Agent works correct but did not save because: "+err.Error(), http.StatusBadRequest)
		        return
            }

            // TODO: Schedule games
            competitors, err := models.GetAgentsN()
            if err != nil{
                panic(err)
            }
            matches := int(math.Ceil(2*math.Log2(float64(competitors))))
            go compete.ScheduleNGames(submission, matches)
           
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
    http.HandleFunc("/api/whoami", getMyName)
    http.HandleFunc("/api/leaderboard", getLeaderboard)
    http.HandleFunc("/api/submissions", getSubmissions)
    http.HandleFunc("/api/matches", getMatches)
    http.HandleFunc("/api/matches/recording", getRecording)
    http.ListenAndServe(":5000", nil)
}
