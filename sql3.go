// https://siongui.github.io/2016/01/09/go-sqlite-example-basic-usage/
// https://github.com/mattn/go-sqlite3/issues/142


package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
        "fmt"
        "time"
        "os"
        "strconv"
)

type TestItem struct {
	Id	string
	Name	string
	Phone	string
	Ph1	string
	Ph2	string
}

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		Id TEXT NOT NULL PRIMARY KEY,
		Name TEXT,
		Phone TEXT,
		Ph1 TEXT,
		Ph2 TEXT,
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
		Name,
		Phone,
		Ph1,
		Ph2,
		InsertedDatetime
	) values(?, ?, ?, ?,?,  CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil { panic(err) }
	defer stmt.Close()

	for _, item := range items {
		_, err2 := stmt.Exec(item.Id, item.Name, item.Phone, item.Ph1, item.Ph2)
		if err2 != nil { panic(err2) }
	}
}

func ReadItem(db *sql.DB) []TestItem {
	sql_readall := `
	SELECT Id, Name, Phone, Ph1, Ph2 FROM items
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil { panic(err) }
	defer rows.Close()

	var result []TestItem
	for rows.Next() {
		item := TestItem{}
		err2 := rows.Scan(&item.Id, &item.Name, &item.Phone, &item.Ph1, &item.Ph2)
		if err2 != nil { panic(err2) }
		result = append(result, item)
	}
	return result
}


type product struct{
    id int
    name string
    phone string
    ph1 string
    ph2 string
    dt time.Time
    //price int
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
        err := rows.Scan(&p.id, &p.name, &p.phone, &p.ph1, &p.ph2,  &p.dt)
        if err != nil{
            fmt.Println(err)
            continue
        }
        products = append(products, p)
    }
    for _, p := range products{
        fmt.Println(p.id, p.name, p.phone, p.dt )
    }
}



func main() {
	const dbpath = "foo.db"
        os.Remove( dbpath )
	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	items := []TestItem{
		TestItem{"1", "A", "213","11","22"},
		TestItem{"2", "B", "214", "33","44"},
	}
	StoreItem(db, items)

	readItems := ReadItem(db)
	fmt.Println(readItems)


//i1, err := strconv.Atoi(str1)
//    if err == nil {
//        fmt.Println(i1)
//    }

	items2 := []TestItem{
		TestItem{"1", "C", "215", strconv.Itoa(111111111), "66"},
		TestItem{"3", "D", "216", "77", "88"},
	}
	StoreItem(db, items2)

	readItems2 := ReadItem(db)
	fmt.Println(readItems2)

        myprint()
}
