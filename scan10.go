// https://github.com/thbar/golang-playground/blob/master/download-files.go
// https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
// https://golang.org/pkg/net/http/


//https://blog.narenarya.in/concurrent-http-in-go.html
//https://medium.com/@dhanushgopinath/concurrent-http-downloads-using-go-32fecfa1ed27



// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779

package main

import (
	"fmt"
	"io"
	"errors"
        "io/ioutil"
	"net/http"
        "log"
	"os"
        "time"
        "runtime"
	"strings"
        "crypto/tls"
        "encoding/xml"
        "math/rand"
        "bytes"
        "bufio"
        "database/sql"
        _ "github.com/mattn/go-sqlite3"


)

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
        if err != nil { panic(err) }
        if db == nil { panic("db nil") }
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
                Id TEXT  NOT NULL PRIMARY KEY,
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
        if err != nil { panic(err) }
}

func StoreItem(items []TestItem) {
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
//db.SetMaxOpenConns(1)
//db.SetMaxIdleConns(1)
// MaxIdleConns <= MaxOpenConns

db := InitDB()
defer db.Close()

 tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

           stmt, err := tx.Prepare(sql_additem)
           if err != nil { panic(err) }
           defer stmt.Close()


        for _, item := range items {
               // https://legacy.gitbook.com/download/pdf/book/astaxie/build-web-application-with-golang?lang=en


               _ , err := tx.Stmt(stmt).Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
               // rows, err := stmt.Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
               //_ , err := tx.Stmt(stmt).Exec(item.Id, item.XmlUrl, item.UrlInfo, item.Ph1, item.Ph2, item.Ph3 )
               if err != nil { panic(err) }
               //tx.Stmt(stmt).Close()


            tx.Commit()
            db.Exec(`PRAGMA shrink_memory;`)

            tx, err = db.Begin()
            if err != nil {
                log.Fatal(err)
            }

 

        }


  tx.Commit()
  db.Exec(`PRAGMA shrink_memory;`)

}




// --------------------------------------------------------------------------------------------------------

 type Item struct {
        XMLName xml.Name `xml:"item"`
        From string `xml:"from"`
        To string `xml:"to"`
        In string `xml:"in"`
        Out string `xml:"out"`
        Amount string `xml:"amount"`
        Param string `xml:"param"`
        Minamount string `xml:"minamount"`
        Maxamount string `xml:"maxamount"`

 }

 type Rates struct {
        XMLName xml.Name `xml:"rates"`
        Items []Item `xml:"item"`
 }

 func (s Item) String() string {
         return fmt.Sprintf("from:%s, to:%s, in:%s, out:%s, amount:%s, param:%s, minamount:%s, maxamount:%s", 
                              s.From, s.To , s.In,  s.Out,  s.Amount,  s.Param,  s.Minamount,  s.Maxamount )
 }


//https://golangbot.com/write-files/

func putLines( d  []Item, path string  ) {  
    f, err := os.Create(path)
    if err != nil {
        fmt.Println(err)
                f.Close()
        return
    }

    for _, v := range d {
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


func mycheck(e error) {
    if e != nil {
        panic(e)
    }
}


//RandomString - Generate a random string of A-Z chars with len = l

func RandomString(len int) string {

      bytes := make([]byte, len)

     for i := 0; i < len; i++ {

          bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25

      }

      return string(bytes)

}


 func mainX( fnm string,  XMLdata []byte   ) (int, error)   {

         //XMLdata, _ := ioutil.ReadAll(xmlFile)
        s1 := RandomString (5)
        items2 := []TestItem{
                TestItem{s1, "C", "215", 111111111, 66, 0},
                //TestItem{"1", "C", "215", 111111111, 66, 0},
            }

       StoreItem(items2)

         var c Rates
         err_xml := xml.Unmarshal(XMLdata, &c)

         if err_xml != nil {
		//fmt.Printf("xml error: ", err_xml)

                 //err := ioutil.WriteFile("bad" + fnm, XMLdata, 0644)
                 //mycheck(err)

                return -1, errors.New("bad xml")
	//	return
         }

         if len (c.Items ) == 0 {
                return -1, errors.New("null len")
         } 

         //fmt.Println(c.Items)
        putLines( c.Items , "ok" + fnm  )


    //for _, num := range c.Items  {
    //    fmt.Println(num) 
   // }

   return 1, nil 
 }





func downloadFromUrl(url string, sec_wait int,  ch chan<-string) {
        start := time.Now()
	tokens := strings.Split(url, "/")
	//fileName := tokens[len(tokens)-1]
	fileName :=  "Y" +  strings.Join(tokens[1:], "") 
	out_xml_data :=  strings.Join(tokens[1:], "X") 

        timeout := time.Duration(time.Duration (sec_wait) * time.Second) 
        //timeout := time.Duration(50  * time.Millisecond) 
        tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify : true}, }
        client := &http.Client{Transport: tr, Timeout: timeout, }
	response, err := client.Get(url)
	//response, err := http.Get(url)
// timeout
	if err != nil {

                if strings.Contains(err.Error(), "Client.Timeout")  {
	               ch <- fmt.Sprintf("%s, Error timeout=%s, url: %s", myt() ,timeout,   url, )

                } else {
                       ch <- fmt.Sprintf("%s, Error url: %s", myt(), url)
                }
		return
	}
	defer response.Body.Close()

        if response.StatusCode != http.StatusOK {
	    ch <- fmt.Sprintf("%s, Error nah=%d, url: %s", myt() , response.StatusCode,   url, )
            return 
        }


	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
	        ch <- fmt.Sprintf("%s, Error creating %s", myt() , fileName, )
		return
	}
	defer output.Close()
// 
        bodyBytes, _ := ioutil.ReadAll(response.Body)
        myerr, err_mess := mainX( out_xml_data,  bodyBytes   )
	if myerr == -1  {
                ch <- fmt.Sprintf("%s, Error %s: %s , file:%s", myt(),   err_mess, url, fileName )
		return
	}
       
        response.Body= ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) 
	n, err := io.Copy(output, response.Body)
	if err != nil {
                ch <- fmt.Sprintf("%s, Error downloading %s-%s", myt(),  url, err)
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
loc, _ := time.LoadLocation("Europe/Moscow")

//set timezone,  
now := time.Now().In(loc)
    return   now.Format("2006-01-02 15:04:05.00")
}



func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}

func isEnable() bool  {
    start, _ := time.Parse(time.RFC822, "20 Feb 14 10:00 UTC")
    end, _ := time.Parse(time.RFC822, "01 Apr 19 10:00 UTC")
    return  inTimeSpan(start, end, time.Now()   )
}


func my1task(sec_wait int ) {

      CreateTable()
      if isEnable() == false {
        fmt.Println( "bad time")
        return
        } 

    arg_len :=len(os.Args)
    f, err := os.OpenFile( os.Args[0] + ".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
      

    //    start := time.Now()
        ch := make(chan string)
	xmlurls := []string{ }




    url_file := "url.txt"
    mylines, err := readLines(url_file)

        if err != nil {
                fmt.Println("Error load %s %s  ", url_file  , err)
                return
        }

        for i := 0; i < len(mylines); i++ {
                raw_url := mylines[i]
                url := strings.TrimSpace(raw_url)
                if  len(url) == 0 {continue}
 
                xmlurls = append(xmlurls, url)
        }


        os.Create("BUSYFLAG")
	for i := 0; i < len(xmlurls); i++ {
		url := xmlurls[i] 
                //fmt.Println(  url)
		go downloadFromUrl(url, sec_wait , ch )
	}
        
        if _, err := f.Write([]byte("//-begin " +  myt() + ", created BUSYFLAG" + "\n"  )); err != nil {
                 log.Fatal(err)
        }

        for range xmlurls{
            msg := <-ch
            if arg_len != 1 {
                   fmt.Println(msg)
            }
            msg += "\n"
            if _, err := f.Write([]byte( msg )); err != nil {
                 log.Fatal(err)
            }
        }
      os.Remove("BUSYFLAG")
  //secs := time.Since(start).Seconds()
  //fmt.Printf("---- time for all ----\n")
  //fmt.Printf("%.2fs secs\n", secs)

//   var input string
//   fmt.Scanln(&input)

        if _, err := f.Write([]byte( "//-end " + myt() + ", removed BUSYFLAG"  +"\n"  )); err != nil {
                 log.Fatal(err)
        }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
 // fmt.Println("scaner: " + myt() )
  return
}



var (
	INTERVAL_SEC = 10 
)

func PrintRoutine1(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
                go my1task(4 )
                time.Sleep(time.Second * 5)
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
        os.Remove("scan9.log")
	runtime.GOMAXPROCS(runtime.NumCPU())
	go PrintRoutine1(INTERVAL_SEC)
	//go PrintRoutine2(INTERVAL_SEC)

	// block forever so that your program won't end
	select {}
}

