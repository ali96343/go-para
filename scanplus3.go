// https://github.com/thbar/golang-playground/blob/master/download-files.go
// https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
// https://golang.org/pkg/net/http/

//https://blog.narenarya.in/concurrent-http-in-go.html
//https://medium.com/@dhanushgopinath/concurrent-http-downloads-using-go-32fecfa1ed27

// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
// https://matt.aimonetti.net/posts/2012/11/27/real-life-concurrency-in-go/

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	// https://habr.com/ru/post/338718/
)

// --------------------------------------------------------------------------------------------------------

// https://habr.com/ru/post/271789/

type DataStore struct {
	sync.Mutex // ← этот мьютекс защищает кэш ниже
	cache      map[string]string
}

func New() *DataStore {
	return &DataStore{
		cache: make(map[string]string),
	}
}

func (ds *DataStore) set(key string, value string) {
	ds.cache[key] = value
}

func (ds *DataStore) get(key string) string {
	if ds.count() > 0 {
		item := ds.cache[key]
		return item
	}
	return ""
}

func (ds *DataStore) count() int {
	return len(ds.cache)
}

func (ds *DataStore) Set(key string, value string) {
	ds.Lock()
	defer ds.Unlock()
	ds.set(key, value)
}

func (ds *DataStore) Get(key string) string {
	ds.Lock()
	defer ds.Unlock()
	return ds.get(key)
}

func (ds *DataStore) Count() int {
	ds.Lock()
	defer ds.Unlock()
	return ds.count()
}

// --------------------------------------------------------------------------------------------------------

var names = []string{"Alan", "Joe", "Jack", "Ben",
	"Ellen", "Lisa", "Carl", "Steve", "Anton", "Yo"}

type my1 struct {
	s1 string
	s2 string
}

type SyncList struct {
	m     sync.Mutex
	slice []interface{}
}

func NewSyncList(cap int) *SyncList {
	return &SyncList{
		sync.Mutex{},
		make([]interface{}, cap),
	}
}

func (l *SyncList) Load(i int) interface{} {
	l.m.Lock()
	defer l.m.Unlock()
	return l.slice[i]
}

func (l *SyncList) Append(val interface{}) {
	l.m.Lock()
	defer l.m.Unlock()
	l.slice = append(l.slice, val)
}

func (l *SyncList) Store(i int, val interface{}) {
	l.m.Lock()
	defer l.m.Unlock()
	l.slice[i] = val
}

var global_pool *SyncList

func test_list() {

	l := NewSyncList(0)
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			l.Append(names[idx])
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < 10; i++ {
		fmt.Printf("Val: %v stored at idx: %d\n", l.Load(i), i)
	}

}

// --------------------------------------------------------------------------------------------------------

type TestItem struct {
	Id      string
	XmlUrl  string
	UrlInfo string
	Ph1     int
	Ph2     int
	Ph3     int
}

var mydb *sql.DB

func InitDB() *sql.DB {
	filepath := "scan.db"
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(1)
	return db
}

func CreateTable() {

	// create table if not exists
	// https://stackoverflow.com/questions/26456253/sqlite-3-not-releasing-memory-in-golang
	// Turn on the Write-Ahead Logging
	sql_table := `
        PRAGMA synchronous = NORMAL;
        PRAGMA temp_store = MEMORY;
        PRAGMA journal_mode = WAL;
        CREATE TABLE IF NOT EXISTS items(
                Id TEXT ,
                XmlUrl TEXT,
                UrlInfo TEXT,
                Ph1 INTEGER,
                Ph2 INTEGER,
                Ph3 INTEGER,
                InsertedDatetime DATETIME
        );
        `
	db := InitDB()
	defer db.Close()
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}

	//db.Exec(`PRAGMA shrink_memory;`)
}

func TruncateTable() {

	// create table if not exists
	// https://stackoverflow.com/questions/26456253/sqlite-3-not-releasing-memory-in-golang
	// Turn on the Write-Ahead Logging
	sql_table := `
        DELETE FROM items;
        VACUUM;
        `
	db := InitDB()
	defer db.Close()
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func pool2string(param interface{}) string {
	s1 := fmt.Sprintf("%s", param)
	s1 = strings.Replace(s1, "{", "", -1)
	return strings.Replace(s1, "}", "", -1)
}

func StoreItem() {
	sql_additem := `
        INSERT INTO items(
                Id,
                XmlUrl,
                UrlInfo,
                Ph1,
                Ph2,
                Ph3,
                InsertedDatetime
        ) values(?, ?, ?, ?,?, ?,  CURRENT_TIMESTAMP)
        `
	//db.SetMaxOpenConns(1)
	//db.SetMaxIdleConns(1)
	// MaxIdleConns <= MaxOpenConns

	var item TestItem

	db := InitDB()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i, e := range global_pool.slice {
		//for i := 0; i < len( global_pool.slice); i++ {
		xx := e.(my1)
		item.Id = strconv.Itoa(i)
		item.XmlUrl = xx.s1
		item.UrlInfo = xx.s2
		item.Ph1 = 0
		item.Ph2 = 0
		item.Ph3 = 0

		// https://legacy.gitbook.com/download/pdf/book/astaxie/build-web-application-with-golang?lang=en

		_, err := tx.Stmt(stmt).Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3)
		// rows, err := stmt.Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
		//_ , err := tx.Stmt(stmt).Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
		if err != nil {
			panic(err)
		}
		//tx.Stmt(stmt).Close()
	}

	tx.Commit()
	db.Exec(`PRAGMA shrink_memory;`)

}

// --------------------------------------------------------------------------------------------------------

type Item struct {
	XMLName   xml.Name `xml:"item"`
	From      string   `xml:"from"`
	To        string   `xml:"to"`
	In        string   `xml:"in"`
	Out       string   `xml:"out"`
	Amount    string   `xml:"amount"`
	Param     string   `xml:"param"`
	Minamount string   `xml:"minamount"`
	Maxamount string   `xml:"maxamount"`
}

type Rates struct {
	XMLName xml.Name `xml:"rates"`
	Items   []Item   `xml:"item"`
}

func (s Item) String() string {
	return fmt.Sprintf("from:%s, to:%s, in:%s, out:%s, amount:%s, param:%s, minamount:%s, maxamount:%s",
		s.From, s.To, s.In, s.Out, s.Amount, s.Param, s.Minamount, s.Maxamount)
}

func myString(s Item) string {
	return fmt.Sprintf("from:%s, to:%s, in:%s, out:%s, amount:%s, param:%s, minamount:%s, maxamount:%s",
		s.From, s.To, s.In, s.Out, s.Amount, s.Param, s.Minamount, s.Maxamount)
}

//https://golangbot.com/write-files/

func putLines(d []Item, path string, orig_url string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	var pool_data my1

	for _, v := range d {
		pool_data.s1 = orig_url
		pool_data.s2 = myString(v)
		global_pool.Append(pool_data)
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("file written successfully: ", path )
	return
}

func put2db() {
	for i := 0; i < len(global_pool.slice); i++ {
		fmt.Println(i, " ", global_pool.slice[i])
	}

}

func mycheck(e error) {
	if e != nil {
		panic(e)
	}
}

//RandomString - Generate a random string of A-Z chars with len = l

func RandomString(len int) string {

	bytes := make([]byte, len)

	for i := 0; i < len; i++ {

		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25

	}

	return string(bytes)

}

func mainX(fnm string, XMLdata []byte, orig_url string) (int, error) {

	//XMLdata, _ := ioutil.ReadAll(xmlFile)
	//s1 := RandomString (5)
	//  items2 := []TestItem{}
	//  items2[0] =    TestItem{s1, "C", "215", 111111111, 66, 0}
	// items2[1] =    TestItem{s1, "C", "215", 111111111, 66, 0}
	//    //TestItem{"1", "C", "215", 111111111, 66, 0},
	// }
	// for i :=0; i  <  1000 ; i++ {
	//   s1 := fnm
	//   s2 := RandomString (5 )
	//   items2 = append( items2, TestItem{s1, s2, "215", 111111111, 66, 0} )

	// }
	//fmt.Println( items2 )
	//go StoreItem(items2)

	var c Rates
	err_xml := xml.Unmarshal(XMLdata, &c)

	if err_xml != nil {
		//fmt.Printf("xml error: ", err_xml)

		//err := ioutil.WriteFile("bad" + fnm, XMLdata, 0644)
		//mycheck(err)

		return -1, errors.New("bad xml")
		//	return
	}

	if len(c.Items) == 0 {
		return -1, errors.New("null len")
	}

	//fmt.Println(c.Items)
	putLines(c.Items, "ok"+fnm, orig_url)

	//for _, num := range c.Items  {
	//    fmt.Println(num)
	// }

	return 1, nil
}

func downloadFromUrl(url string, sec_wait int, ch chan<- string) {
	start := time.Now()
	tokens := strings.Split(url, "/")
	//fileName := tokens[len(tokens)-1]
	fileName := "Y" + strings.Join(tokens[1:], "")
	out_xml_data := strings.Join(tokens[1:], "X")

	timeout := time.Duration(time.Duration(sec_wait) * time.Second)
	//timeout := time.Duration(50  * time.Millisecond)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr, Timeout: timeout}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- fmt.Sprintf("%s, Error reques url: %s", myt(), url)
		return
	}

	req.Header.Set("User-Agent", "My User Agent 2.0")
	//req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	response, err := client.Do(req)

	//response, err := client.Get(url)
	//response, err := http.Get(url)

	if response != nil {
		defer response.Body.Close()
	}

	// timeout
	if err != nil {

		if strings.Contains(err.Error(), "Client.Timeout") {
			ch <- fmt.Sprintf("%s, Error timeout=%s, url: %s", myt(), timeout, url)

		} else {
			ch <- fmt.Sprintf("%s, Error (unk) url: %s", myt(), url)
		}
		return
	}

	if response.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("%s, Error nah=%d, url: %s", myt(), response.StatusCode, url)
		return
	}

	if response.ContentLength == 0 {
		ch <- fmt.Sprintf("%s, Error nahnil, url: %s", myt(), url)
		return
	}

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		ch <- fmt.Sprintf("%s, Error creating %s", myt(), fileName)
		return
	}
	defer output.Close()
	//
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	myerr, err_mess := mainX(out_xml_data, bodyBytes, url)
	if myerr == -1 {
		ch <- fmt.Sprintf("%s, Error %s: %s , file:%s", myt(), err_mess, url, fileName)
		return
	}

	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	n, err := io.Copy(output, response.Body)
	if err != nil {
		ch <- fmt.Sprintf("%s, Error downloading %s-%s", myt(), url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%s, %.2f sec, len: %d, url: %s, csv: ok%s", myt(), secs, n, url, out_xml_data)
	return
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func myt() string {
	loc, _ := time.LoadLocation("UTC")

	//set timezone,
	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05.00")
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func isEnable() bool {
	start, _ := time.Parse(time.RFC822, "20 Feb 14 10:00 UTC")
	end, _ := time.Parse(time.RFC822, "25 Apr 19 10:00 UTC")
	return inTimeSpan(start, end, time.Now())
}

func my1task(sec_wait int) {
	if isEnable() == false {
		fmt.Println("bad time")
		return
	}

	arg_len := len(os.Args)
	f, err := os.OpenFile(os.Args[0]+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//    start := time.Now()
	ch := make(chan string)
	defer close(ch)
	xmlurls := []string{}

	url_file := "url.txt"
	mylines, err := readLines(url_file)

	if err != nil {
		fmt.Println("Error load %s %s  ", url_file, err)
		return
	}

	for i := 0; i < len(mylines); i++ {
		raw_url := mylines[i]
		url := strings.TrimSpace(raw_url)
		if len(url) == 0 {
			continue
		}

		xmlurls = append(xmlurls, url)
	}

	global_pool = NewSyncList(0)

	os.Create("BUSYFLAG")
	defer os.Remove("BUSYFLAG")
	for i := 0; i < len(xmlurls); i++ {
		url := xmlurls[i]
		//fmt.Println(  url)
		go downloadFromUrl(url, sec_wait, ch)
	}

	if _, err := f.Write([]byte("//-begin " + myt() + ", set BUSYFLAG" + "\n")); err != nil {
		log.Fatal(err)
	}

	for range xmlurls {
		msg := <-ch
		if arg_len != 1 {
			fmt.Println(msg)
		}
		msg += "\n"
		if _, err := f.Write([]byte(msg)); err != nil {
			log.Fatal(err)
		}
	}
	//secs := time.Since(start).Seconds()
	//fmt.Printf("---- time for all ----\n")
	//fmt.Printf("%.2fs secs\n", secs)

	//   var input string
	//   fmt.Scanln(&input)

	// fmt.Println("scaner: " + myt() )
	f.Write([]byte("//end " + myt() + ", rm BUSYFLAG" + "\n"))

	f.Write([]byte("//lock db for replace all data " + myt() + ", set LOCKDBFLAG" + "\n"))
	os.Create("LOCKDBFLAG")
	defer os.Remove("LOCKDBFLAG")
	//fmt.Println(global_pool )
	//put2db()
	TruncateTable()
	StoreItem()
	//global_pool = NewSyncList(0)
	global_pool = nil
	f.Write([]byte("//unlock db: data updated " + myt() + ", rm  LOCKDBFLAG" + "\n"))

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return
}

var (
	INTERVAL_SEC = 8
	//INTERVAL_SEC = 10
)

func PrintRoutine1(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
		go my1task(4)
		time.Sleep(time.Second * 4)
		//fmt.Println("from t1")
	}
}

func PrintRoutine2(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
		fmt.Println("from t2")
	}
}

func main() {
	// https://socketloop.com/tutorials/golang-call-a-function-after-some-delay-time-sleep-and-tick

	os.Remove("scan.db")
	CreateTable()
	runtime.GOMAXPROCS(runtime.NumCPU())
	go PrintRoutine1(INTERVAL_SEC)
	//go PrintRoutine2(INTERVAL_SEC)

	// block forever so that your program won't end
	select {}
}
