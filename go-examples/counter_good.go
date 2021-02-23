package main

import (
	"fmt"
	//"time"
    "sync"
)

func counter(count *int, numTimes int, wg *sync.WaitGroup, mu *sync.Mutex) {
    defer wg.Done()
	for i := 0; i < numTimes; i++ {
        mu.Lock()
		(*count) = (*count) + 1
        mu.Unlock()
	}
}

func main() {
    count := 0
	var wait sync.WaitGroup
    var lock sync.Mutex
    wait.Add(2)
	go counter(&count,100000,&wait,&lock)
    go counter(&count,100000,&wait,&lock)
    wait.Wait()
    fmt.Println("After incrementing 200000 times, count is",count)
}
