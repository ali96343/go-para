package main

import (
	"fmt"
	"sync"
)

type my1 struct {
   s1 string
   s2 string
}


var names = []string{"Alan", "Joe", "Jack", "Ben", "Ellen", "Lisa", "Carl", "Steve", "Anton", "Yo"}
var n1 = []my1{  {"Alan", "Joe"}, {"Jack", "Ben"}, { "Ellen", "Lisa"}, {"Carl", "Steve"}, {"Anton", "Yo"}  }

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
// https://github.com/dnaeon/gru/blob/master/utils/slice.go

func teste(param interface{}) string {
        return fmt.Sprintf("%s", param)
//        return *param.(string) 
}



func main() {
global_pool = NewSyncList(0)
	l := NewSyncList(0)
	wg := &sync.WaitGroup{}
	wg.Add(len(n1))
	for i := 0; i < len(n1); i++ {
		go func(idx int) {
			l.Append(n1[idx])
			//l.Append(names[idx])
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := 0; i < len(n1); i++ {
		fmt.Printf("Val: %v stored at idx: %d\n", l.Load(i), i)
	}
          
        
fmt.Println(l )





for i , e  := range(l.slice) {
 fmt.Println(i,e)
fmt.Printf("%T\n", e)

  x := e.(my1)
 fmt.Println(i, x.s1  )
 fmt.Println(i, x.s2  )

 fmt.Println(i, e  )
}



fmt.Println(len(l.slice ))
fmt.Printf("%T\n", l)
global_pool = nil
}
