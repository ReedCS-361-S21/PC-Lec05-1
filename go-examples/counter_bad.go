package main

import (
	"fmt"
	//"time"
    "sync"
)

func counter(count *int, numTimes int, wg *sync.WaitGroup) {
    defer wg.Done()
	for i := 0; i < numTimes; i++ {
		(*count) = (*count) + 1
	}
}

func main() {
    count := 0
	var wait sync.WaitGroup
    wait.Add(2)
	go counter(&count,100000,&wait)
    go counter(&count,100000,&wait)
    wait.Wait()
    fmt.Println("After incrementing 200000 times, count is",count)
}
