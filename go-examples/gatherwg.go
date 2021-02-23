package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func fillup(I int, xs []int, N int, P int, wait *sync.WaitGroup) {
	defer wait.Done()
	n := N / P
	for i := 0; i < n; i++ {
		xs[I*n+i] = rand.Intn(N)
	}
}

func sumup(I int, ss[] int, xs []int, N int, P int, wait *sync.WaitGroup) {
	defer wait.Done()
	n := N / P
	s := 0
	for i := 0; i < n; i++ {
		s += xs[I*n+i]
	}
	ss[I] = s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	P := 4
	N := 256
	if len(os.Args) > 1 {
		sz, err := strconv.Atoi(os.Args[1])
		if err == nil {
			N = sz
		}
	}
	if len(os.Args) > 2 {
		ps, err := strconv.Atoi(os.Args[2])
		if err == nil {
			P = ps
		}
	}

	var values []int = make([]int, N)
	results := make([]int, P)

	// Initialize the data (with several threads).
	var wait1 sync.WaitGroup
	for I := 0; I < P; I++ {
		go fillup(I, values, N, P, &wait1)
	}
	wait1.Wait()

	// Sum up the data in parallel, timing the parallel work.
	var wait2 sync.WaitGroup
	wait2.Add(P)
	start := time.Now()
	for I := 0; I < P; I++ {
		go sumup(I, results, values, N, P, &wait2)
	}
	wait2.Wait()
	
	duration := time.Since(start)
  	sum := 0
        for I := 0; I < P; I++ {
   	        sum += results[I]
        }
   	fmt.Println("The sum computed was",sum)
	fmt.Println("Time to compute:", duration)
	fmt.Println(P,"workers ran on",runtime.GOMAXPROCS(0),"processors.")
	fmt.Println()
}
