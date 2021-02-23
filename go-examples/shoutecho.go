package main

import (
	"fmt"
	"math/rand"
	"time"
)

func shout(c chan int) {
	var count int = 0
	for {
		c <- count
		count = count + 1
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func echo(c chan int) {
	for {
		var count int = <-c
		fmt.Println(count)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func main() {
	var c chan int = make(chan int)
	go shout(c)
	go echo(c)
	for {
		time.Sleep(time.Second)
	}
}
