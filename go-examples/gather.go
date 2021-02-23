package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func randomPerm(size int) []int {
	// Generate a random permutation.
	xs := make([]int, size)
	for i := 0; i < size; i += 1 {
		xs[i] = i
	}
	for i := 0; i < size; i += 1 {
		j := rand.Intn(size-i) + i
		xs[i], xs[j] = xs[j], xs[i]
	}
	return xs
}

func sumup(I int, xs []int, N int, P int, r chan int) {
	n := N / P
	s := 0
	for i := 0; i < n; i++ {
		s += xs[I*n+i]
	}
	r <- s
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
	var values []int = randomPerm(N)
	results := make(chan int)
	start := time.Now()
	for I := 0; I < P; I++ {
		go sumup(I, values, N, P, results)
	}
	sum := 0
	for I := 0; I < P; I++ {
		sum += <-results
	}
	duration := time.Since(start)
	fmt.Println("The sum computed was", sum)
	fmt.Println("The sum expected was", N*(N-1)/2)
	fmt.Println("Time to compute:", duration)
}
