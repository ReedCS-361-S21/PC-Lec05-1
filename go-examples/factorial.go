package main

import "fmt"

func factorial(n int) int {
	p := 1
	for i := 1; i < n; i++ {
		p *= i
	}
	return p
}

func main() {
	fmt.Println(factorial(10))
}
