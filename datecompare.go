package main

import (
    "fmt"
    "time"
)

func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}

func isEnable() bool  {
    start, _ := time.Parse(time.RFC822, "20 Feb 14 10:00 UTC")
    end, _ := time.Parse(time.RFC822, "01 Apr 18 10:00 UTC")

    in := time.Now()

    return  inTimeSpan(start, end, in) 

}



func main () {
      if isEnable() == true {
        fmt.Println( "ok")
        } else {
        fmt.Println( "bad")

       }

    }
