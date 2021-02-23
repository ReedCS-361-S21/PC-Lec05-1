package main

import (
	"fmt"
	"time"
)

func walk() {
	for {
		fmt.Printf("step\n")
	}
}

func chewGum() {
	for {
		fmt.Printf("chomp\n")
	}
}

func main() {
	fmt.Printf("Starting...\n")
	go walk()
	go chewGum()
	for {
		time.Sleep(time.Second)
	}
}
