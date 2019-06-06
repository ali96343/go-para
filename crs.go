package main

import (
    "fmt"
)


type xx interface {}


func print_empty( X xx ) {
     ss := X.(type) 
     fmt.Printf( "  %#v \n", X   )
     fmt.Printf( "  %#T \n", X   )
     s:= X.(string)
     fmt.Println(s)

}

func main() {
   fmt.Println("hi")
   print_empty( "1111111111" )
}
