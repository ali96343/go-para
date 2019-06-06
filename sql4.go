// https://siongui.github.io/2016/01/09/go-sqlite-example-basic-usage/
// https://github.com/mattn/go-sqlite3/issues/142
// https://metanit.com/go/tutorial/10.4.php
// https://stackoverflow.com/questions/35804884/sqlite-concurrent-writing-performance
//https://www.alexedwards.net/blog/configuring-sqldb

package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
        "fmt"
        "time"
        "os"
        //"strconv"
)

type TestItem struct {
	Id	int
	XmlUrl	string
	UrlInfo	string
	Ph1     int	
	Ph2     int	
	Ph3     int	
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func CreateTable(db *sql.DB) {

	// create table if not exists
        // https://stackoverflow.com/questions/26456253/sqlite-3-not-releasing-memory-in-golang

	sql_table := `
        PRAGMA journal_mode = WAL;
	CREATE TABLE IF NOT EXISTS items(
		Id INTEGER  NOT NULL PRIMARY KEY,
		XmlUrl TEXT,
		UrlInfo TEXT,
		Ph1 INTEGER,
		Ph2 INTEGER,
		Ph3 INTEGER,
		InsertedDatetime DATETIME
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }
}

func StoreItem(db *sql.DB, items []TestItem) {
	sql_additem := `
	INSERT OR REPLACE INTO items(
		Id,
		XmlUrl,
		UrlInfo,
		Ph1,
		Ph2,
		Ph3,
		InsertedDatetime
	) values(?, ?, ?, ?,?, ?,  CURRENT_TIMESTAMP)
	`
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
// MaxIdleConns <= MaxOpenConns

	stmt, err := db.Prepare(sql_additem)
	if err != nil { panic(err) }
	defer stmt.Close()

       
	for _, item := range items {

tx, err := db.Begin()
if err != nil {
fmt.Println(err)
}
// https://legacy.gitbook.com/download/pdf/book/astaxie/build-web-application-with-golang?lang=en
		_, err2 := stmt.Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
	//	if err2 != nil { panic(err2) }

if err2 != nil {
fmt.Println("doing rollback")
tx.Rollback()
} else {
tx.Commit()
}


	}




}

func ReadItem(db *sql.DB) []TestItem {
	sql_readall := `
	SELECT Id, XmlUrl, UrlInfo, Ph1, Ph2, Ph3 FROM items
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil { panic(err) }
	defer rows.Close()

	var result []TestItem
	for rows.Next() {
		item := TestItem{}
		err2 := rows.Scan(&item.Id, &item.XmlUrl, &item.UrlInfo, &item.Ph1, &item.Ph2, &item.Ph3 )
		if err2 != nil { panic(err2) }
		result = append(result, item)
	}
	return result
}


type product struct{
    id int
    xmlUrl string
    url_info string
    ph1 int
    ph2 int
    ph3 int
    dt time.Time
}
func myprint() { 
    //os.Remove("foo.db" ) 
    db, err := sql.Open("sqlite3", "foo.db")
    if err != nil {
        fmt.Println("panic")
        panic(err)
    }
    defer db.Close()
    rows, err := db.Query("select * from items")
    if err != nil {
        fmt.Println("panic")
        panic(err)
    }
    defer rows.Close()
    products := []product{}
     
    for rows.Next(){
        p := product{}
        err := rows.Scan(&p.id, &p.xmlUrl, &p.url_info, &p.ph1, &p.ph2,  &p.ph3, &p.dt)
        if err != nil{
            fmt.Println(err)
            continue
        }
        products = append(products, p)
    }
    for _, p := range products{
        fmt.Println(" ", p.id, " ", p.xmlUrl, " ", p.url_info, " ", p.ph1, " ", p.ph2, " ", p.ph3, " ",  p.dt, )
    }
}

func findurl ( fnd_id int ) {

sqlStatement := `SELECT * FROM items WHERE id=$1;`
var p product 

    db, err := sql.Open("sqlite3", "foo.db")
    if err != nil {
        fmt.Println("panic")
        panic(err)
    }
    defer db.Close()


row := db.QueryRow(sqlStatement,  fnd_id  )
err = row.Scan(&p.id, &p.xmlUrl, &p.url_info, &p.ph1, &p.ph2,  &p.ph3, &p.dt)

switch err {
case sql.ErrNoRows:
  fmt.Println("No rows were returned!")
  return
case nil:
        fmt.Println(" ", p.id, " ", p.xmlUrl, " ", p.url_info, " ", p.ph1, " ", p.ph2, " ", p.ph3, " ",  p.dt, )
default:
  panic(err)
}

}


func main() {
	const dbpath = "foo.db"
        os.Remove( dbpath )
	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	items := []TestItem{
		TestItem{1, "A", "213",11,22, 0},
		TestItem{2, "B", "214", 33,44, 0},
	}
	StoreItem(db, items)

	readItems := ReadItem(db)
	fmt.Println(readItems)


//i1, err := strconv.Atoi(str1)
//    if err == nil {
//        fmt.Println(i1)
//    }

	items2 := []TestItem{
		TestItem{1, "C", "215", 111111111, 66, 0},
		TestItem{3, "D", "216", 77, 88, 0},
		TestItem{4, "D", "216", 77, 88, 0},
		TestItem{5, "D", "216", 77, 88, 0},
	}
	StoreItem(db, items2)

	readItems2 := ReadItem(db)
	fmt.Println(readItems2)

        myprint()

        findurl ( 1  )
}
