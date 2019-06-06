package main
import (

   "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// https://metanit.com/go/tutorial/10.4.php
// https://siongui.github.io/2016/01/09/go-sqlite-example-basic-usage/

func main() { 
 
    db, err := sql.Open("sqlite3", "adb/test.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    result, err := db.Exec("insert into cars  (name, price) values ( $1, $2)", 
        "ZAPOR", 72000)
    if err != nil{
        panic(err)
    }
    fmt.Println(result.LastInsertId())  // id последнего добавленного объекта
    fmt.Println(result.RowsAffected())  // количество добавленных строк
     
}
