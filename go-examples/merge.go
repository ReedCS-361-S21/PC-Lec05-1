package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func merge(A []int, B []int) []int {
	var C []int = make([]int, len(A)+len(B))
	i := 0
	j := 0
	for k := 0; k < len(C); k++ {
		if j == len(B) {
			C[k] = A[i]
			i = i + 1
		} else if i == len(A) {
			C[k] = B[j]
			j = j + 1
		} else if A[i] <= B[j] {
			C[k] = A[i]
			i = i + 1
		} else {
			C[k] = B[j]
			j = j + 1
		}
	}
	return C
}

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

func split(A []int) ([]int, []int) {
	n := len(A)
	L := make([]int, n-n/2)
	R := make([]int, n/2)
	for i := 0; i < len(L); i++ {
		L[i] = A[i]
	}
	for i := 0; i < len(R); i++ {
		R[i] = A[len(L)+i]
	}
	return L, R
}

func bubble(A []int) {
	for i := 1; i < len(A); i++ {
		for j := i; j > 0; j-- {
			if A[j-1] > A[j] {
				A[j], A[j-1] = A[j-1], A[j]
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	size := 20
	if len(os.Args) > 1 {
		s, err := strconv.Atoi(os.Args[1])
		if err == nil {
			size = s
		}
	}
	var values []int = randomPerm(size)
	var a []int
	var b []int
	a, b = split(values)
	bubble(a)
	bubble(b)
	fmt.Println(a)
	fmt.Println(b)
	var c []int = merge(a, b)
	fmt.Println(c)
}
