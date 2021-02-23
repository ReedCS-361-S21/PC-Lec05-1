package main

import (
	"fmt"
	"strings"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func flipcap(s string, p int) string {
	if "a" <= s[p:p+1] && s[p:p+1] <= "z" {
		return s[:p] + strings.ToUpper(s[p:p+1]) + s[p+1:]
	} else if "A" <= s[p:p+1] && s[p:p+1] <= "Z" {
		return s[:p] + strings.ToLower(s[p:p+1]) + s[p+1:]
	} else {
		return s
	}
}

func forward(I int, T int, pred chan string, succ chan string) {
	for t := 1; t < T; t++ {
		m := <-pred
		fmt.Printf("%02d. %02d: Received %s.\n", t, I, m)
		m = flipcap(m, rand.Intn(len(m)))
		fmt.Printf("%02d. %02d: Sent %s.\n", t, I, m)
		succ <- m
	}
}

func main() {
	T := 3
	if len(os.Args) > 1 {
		ts, err := strconv.Atoi(os.Args[1])
		if err == nil {
			T = ts
		}
	}
	message := "hello"
	if len(os.Args) > 2 {
		message = os.Args[2]
	}
	P := 4
	if len(os.Args) > 3 {
		ps, err := strconv.Atoi(os.Args[3])
		if err == nil {
			P = ps
		}
	}

	rand.Seed(time.Now().UnixNano())

	chs := make([]chan string, P)
	for i := 0; i < P; i++ {
		chs[i] = make(chan string)
	}

	for i := 1; i < P; i++ {
		go forward(i, T, chs[i], chs[(i+1)%P])
	}

	chs[1] <- message
	forward(0, T-1, chs[0], chs[1])
	message = <-chs[0]
}
