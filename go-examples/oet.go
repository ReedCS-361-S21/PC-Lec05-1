package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type info struct {
	source int
	value  int
}

func main() {
	size := 3
	if len(os.Args) > 1 {
		s, err := strconv.Atoi(os.Args[1])
		if err == nil {
			size = s
		}
	}
	// Create a random permutation of values.
	rand.Seed(time.Now().UnixNano())
	values := randomPerm(size)

	// Create channels to ship values to/from the proc array.
	inputs := make(chan int)
	results := make(chan info)

	// Print unsorted data.
	output(values)

	// Start the sorting processors.
	makeProcArray(size, inputs, results)

	// Scatter the data to them.
	for v := range values {
		inputs <- v
	}

	// Gather the sorted data from them.
	for i := 0; i < size; i++ {
		sd := <-results
		values[sd.source] = sd.value
	}

	// Print sorted data.
	output(values)
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

//
//            | in
//            V
//         +-----+
// outL <--| I   |--> outR
//  inL -->| of P|<-- inR
//         +-----+
//            |
//            V out
//
func oddEvenProc(I int, P int,
	inL, outL chan int,
	inR, outR chan int,
	in chan int, out chan info) {
	data := <-in
	for t := 0; t < P; t += 1 {
		if (I+t)%2 == 0 && I < P-1 {
			outR <- data
			next := <-inR
			if next < data {
				data = next
			}
		}
		if (I+t)%2 == 1 && I > 0 {
			next := <-inL
			outL <- data
			if next > data {
				data = next
			}
		}
	}

	out <- info{source: I, value: data}
}

func makeProcArray(P int, in chan int, out chan info) {
	var inL chan int = nil
	var inR chan int = nil
	var outL chan int = nil
	var outR chan int = nil
	for I := 0; I < P; I++ {
		inL = outR
		outL = inR
		if I != P-1 {
			inR = make(chan int)
			outR = make(chan int)
		} else {
			inR = nil
			outR = nil
		}
		go oddEvenProc(I, P, inL, outL, inR, outR, in, out)
	}
}

func output(values []int) {
	for _, v := range values {
		fmt.Print(v, " ")
	}
	fmt.Println()
}
