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
        "bytes"
        "bufio"
)

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
    //d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}

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
}


func mycheck(e error) {
    if e != nil {
        panic(e)
    }
}


 func mainX( fnm string,  XMLdata []byte   ) (int, error)   {

         //XMLdata, _ := ioutil.ReadAll(xmlFile)

         var c Rates
         err_xml := xml.Unmarshal(XMLdata, &c)

         if err_xml != nil {
		//fmt.Printf("xml error: ", err_xml)

                 err := ioutil.WriteFile("bad" + fnm, XMLdata, 0644)
                 mycheck(err)

                return -1, errors.New("bad xml")
	//	return
         }

         //fmt.Println(c.Items)
         putLines( c.Items , "ok" + fnm  )

    //for _, num := range c.Items  {
    //    fmt.Println(num) 
   // }

   return 1, nil 
 }





func downloadFromUrl(url string, ch chan<-string) {
        start := time.Now()
	tokens := strings.Split(url, "/")
	//fileName := tokens[len(tokens)-1]
	fileName :=  strings.Join(tokens[1:], "") 
	file2 :=  strings.Join(tokens[1:], "X") 

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		//fmt.Println(myt()+  " " + "Error while creating", fileName, "-", err)
	        ch <- fmt.Sprintf("%s, Error creating %s", myt() , fileName, )
		return
	}
	defer output.Close()

      
        timeout := time.Duration(5 * time.Second) 
        tr := &http.Transport{ TLSClientConfig: &tls.Config{InsecureSkipVerify : true}, }
        client := &http.Client{Transport: tr, Timeout: timeout, }
	response, err := client.Get(url)
	//response, err := http.Get(url)

	if err != nil {
	//	fmt.Println("Error while downloading", url, "-", err)

           if strings.Contains(err.Error(), "Client.Timeout")  {
	           ch <- fmt.Sprintf("%s, Error timeout=%s, url: %s", myt() ,timeout,   url, )

                  } else {
                ch <- fmt.Sprintf("%s, Error url %s", myt(), url)
                }
		return
	}

	defer response.Body.Close()

        bodyBytes, _ := ioutil.ReadAll(response.Body)
        myerr, err_mess := mainX( file2,  bodyBytes   )
	if myerr == -1  {
		//fmt.Println("Error while downloading", url, "-", err)
           ch <- fmt.Sprintf("%s, Error %s: %s", myt(),   err_mess, url)
		return
	}
       
        response.Body= ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) 

	n, err := io.Copy(output, response.Body)
        //fmt.Println(n, ": bytes") 
	if err != nil {
		//fmt.Println("Error while downloading", url, "-", err)
           ch <- fmt.Sprintf("%s, Error downloading %s-%s", myt(),  url, err)
		return
	}
        secs := time.Since(start).Seconds()
        //fmt.Printf("%.2fs elapsed\n",  secs  )
	//fmt.Println("---- Downloading", url, "to", fileName)
	//fmt.Println(n, "bytes downloaded.")
        //ch <- fmt.Sprintf("%.2f sec,  %s , fnm= %s", secs,  url, file2)
        ch <- fmt.Sprintf("%s, %.2f sec, len: %d, url: %s, csv: %s", myt(), secs, n, url, file2)
        //ch <- fmt.Sprintf("%.2f sec, length: %d,  %s , fnm= %s", secs, n , url, fileName)
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
    return   time.Now().Format("2006-01-02 15:04:05")
}


func my1task( ) {


    arg_len :=len(os.Args)
    f, err := os.OpenFile( os.Args[0] + ".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

        start := time.Now()
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

	for i := 0; i < len(xmlurls); i++ {
		url := xmlurls[i] 
                //fmt.Println(  url)
		go downloadFromUrl(url, ch )
	}
        
        if _, err := f.Write([]byte("//-begin " +  myt() + "\n"  )); err != nil {
                 log.Fatal(err)
        }

        for range xmlurls{
            msg := <-ch
            if arg_len == 1 {
                   fmt.Println(msg)
            }
            msg += "\n"
            if _, err := f.Write([]byte( msg )); err != nil {
                 log.Fatal(err)
            }
        }

  secs := time.Since(start).Seconds()
  fmt.Printf("---- time for all ----\n")
  fmt.Printf("%.2fs secs\n", secs)

//   var input string
//   fmt.Scanln(&input)

        if _, err := f.Write([]byte( "//-end " + myt() + "\n"  )); err != nil {
                 log.Fatal(err)
        }
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
  fmt.Println("scaner: " + myt() )
  time.Sleep(time.Second * 5)
  return
}



var (
	INTERVAL_SEC = 10 
)

func PrintRoutine1(intervalInSec int) {
	t := time.NewTicker(time.Duration(intervalInSec) * time.Second)
	for _ = range t.C {
                go my1task()
		fmt.Println("from t1")
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

	runtime.GOMAXPROCS(runtime.NumCPU())
	go PrintRoutine1(INTERVAL_SEC)
	//go PrintRoutine2(INTERVAL_SEC)

	// block forever so that your program won't end
	select {}
}

