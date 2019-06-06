package main

import (
	"fmt"
	"sync"
)

// ------------------------------------------------------------------------------------

// Slice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// Concurrent slice item
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// NewConcurrentSlice creates a new concurrent slice
func NewConcurrentSlice() *ConcurrentSlice {
	cs := &ConcurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Unlock()
		for index, value := range cs.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}


// ------------------------------------------------------------------------------------



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

func main() {
	store := New()
	store.Set("Go", "Lang")
	result := store.Get("Go")
	fmt.Println(result)
}

