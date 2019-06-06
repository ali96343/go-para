package main

import (
    "database/sql"
    "fmt"
    "strconv"

    _ "github.com/mattn/go-sqlite3"
)

func main() {
    database, _ := sql.Open("sqlite3", "statscan.db")
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS scanurls (id INTEGER PRIMARY KEY, url TEXT, numget INTEGER, numnah INTEGER)")
    statement.Exec()
    statement, _ = database.Prepare("INSERT INTO scanurls  (url, numget, numnah) VALUES (?, ?, ?)")
    statement.Exec("url", 0, 0)
    rows, _ := database.Query("SELECT id, url, numget, numnah  FROM scanurls")
    var id int
    var url  string
    var numget  string
    var numnah  string
    for rows.Next() {
        rows.Scan(&id, &url, &numget, &numnah)
        fmt.Println(strconv.Itoa(id) + ": " + numget  + " " + numnah)
    }
}
