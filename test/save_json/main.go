package main

import (
    "database/sql"
    "github.com/lib/pq"
    "fmt"
)
func main() {

    var DBurl string = "postgresql://magisterbrownie:post@localhost/subm_ranking?sslmode=disable"
    DB, _ := sql.Open("postgres", DBurl) 
    //_, err := DB.Exec("INSERT INTO matches (subm_1, subm_2, subm_1_change, subm_2_change) VALUES (1, 2, 112, -342)")

    //err := DB.QueryRow("INSERT INTO submissions (user_id, file_name,  container_id, raiting, sigma) VALUES (2, 'assa.tar', 'asa', 23, 32) RETURNING id, created_at")
    _, err := DB.Exec("SELECT * FROM submissions WHERE id != ALL($1)", pq.Array([]int{})) 

    fmt.Println(err)


}
