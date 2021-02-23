package main

import (
	"fmt"
	"math/rand"
	"os"
	//"runtime"
	"strconv"
	"sync"
	"time"
)

type shared struct {
    sum int
    count int
    mu *sync.Mutex
    cv *sync.Cond
}
func fillup(I int, xs []int, N int, P int, wait *sync.WaitGroup) {
	defer wait.Done()
	n := N / P
	for i := 0; i < n; i++ {
		xs[I*n+i] = rand.Intn(N)
	}
}

func sumup(I int, xs []int, N int, P int, sh *shared) {
	n := N / P
	s := 0
	for i := 0; i < n; i++ {
		s += xs[I*n+i]
	}
	sh.mu.Lock()
    sh.sum += s
    sh.count++
    sh.cv.Signal()
    sh.mu.Unlock()
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

	// Initialize the data (with several threads).
	var wait1 sync.WaitGroup
    wait1.Add(P)
	for I := 0; I < P; I++ {
		go fillup(I, values, N, P, &wait1)
	}
	wait1.Wait()

    if N <= 24 {
        fmt.Println(values)
    }

	// Sum up the data in parallel, timing the parallel work.
    var lock sync.Mutex
    condition := sync.NewCond(&lock)
	tally := &shared{sum:0, count:0, mu:&lock, cv:condition}

	for I := 0; I < P; I++ {
		go sumup(I, values, N, P, tally)
	}

    tally.mu.Lock()
    for tally.count < P {
        tally.cv.Wait()
    }
    tally.mu.Unlock()

   	fmt.Println("The sum computed was",tally.sum)
}
