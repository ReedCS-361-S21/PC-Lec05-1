package main

// Reed College CSCI 361
//
// oetsort.go
//
// Odd-Even Transposition Sort in Go.
//
// See the comments for `func main` below, which describe how to run the program.

import (
	"fmt"
	"os"
	"strconv"
    "sync"
    "math/rand"
	"time"
)

// report struct
//
// Used by a "processor" to report its value to the "front end".
type report struct {
	source int
	value  int
}

// procinfo struct
//
// Info pertaining to a "processor"'s configuration within its "linear array".
type procinfo struct {
    id int
    procs int
    inLeft, outLeft, inRight, outRight chan int
    input chan int
    result chan report
    signal *sync.WaitGroup
}

// Usage:
//     go run oetsort.go [<size=3>]
//
// Runs a simulation of Odd-Even-Transposition-Sort using <size> data values
// and mimicking a linear processor array of the same size.
//
// This uses a goroutine worker for each processor and communicates using Go
// channels.
//
// A single run of OET is applied to a random permutation of [0..size-1].
//
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
	inputs := make(chan int, size)
	results := make(chan report, size)

	// Print unsorted data.
	output2d(values)

    var done sync.WaitGroup
    done.Add(size);

	// Create the processor array to perform oetSort.
	makeProcArray(size, inputs, results, oetSort, &done)

	// Scatter the data to them.
	for v := range values {
		inputs <- v
	}

    done.Wait()

	// Gather the sorted data from them.
	for i := 0; i < size; i++ {
		sd := <-results
		values[sd.source] = sd.value
	}

	// Print sorted data.
	output2d(values)
}

// randomPerm(size):
//
// Build and return an array that's a permutation of [0..size-1].
//
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

// mini(x,y):
//
// Return the minimum of two integers.
//
func mini(x int, y int) int {
    // Return the minimum of two integers.
    if x <= y {
        return x
    } else {
        return y
    }
}

// maxi(x,y):
//
// Return the maximum of two integers.
//
func maxi(x int, y int) int {
    if x >= y {
        return x
    } else {
        return y
    }
}

// exchangeWith(op, value, in, out):
//
// Send your `value` to another processor over `out`.
// Receive their value over `in`.
// Combine yours with theirs using `op`, updating yours.
//
func exchangeWith(op func(int,int)int, value *int, in chan int, out chan int) {
    mine := (*value)
    out <- mine
    other := <- in
    (*value) = op(mine,other)
}

// oetRun(value, procid, steps, toPred, fromPred, toSucc, fromSucc):
//
// Make `steps`  of Odd-Even-Transposition-Sort as `procid`.
// Update integer in `*value`.
//
func oetRun(value *int, procid int, steps int,
            toPred, fromPred, toSucc, fromSucc chan int) {
    for t := 0; t < steps; t += 1 {
        if (procid + t) % 2 == 0 && toSucc != nil {
            exchangeWith(mini, value, fromSucc, toSucc)
        }
        if (procid + t) % 2 == 1 && toPred != nil {
            exchangeWith(maxi, value, fromPred, toPred)
        }
    }
}

// oetSort(procinfo):
//
// Participate in a full run of Odd-Even-Tramnsposition-Sort as a
// processing element in a one-dimensional array.
//
// Your connections are given by `proc` as shown below.
//
//                | input
//                V
//            +-------+
// outLeft <--| value |--> outRight
//  inLeft -->|       |<-- inRight
//            +-------+
//                |
//                V result
//
// You get your initial value from `input` channel.
// You coordinate with the array using the `in*` and `out*` channels.
// You report your final value using the `result` channel.
// You then signal that you are done before exiting.
//
func oetSort(proc procinfo) {
    defer proc.signal.Done()
	value := <-proc.input
    oetRun(&value, proc.id, proc.procs,
           proc.outLeft, proc.inLeft,
           proc.outRight, proc.inRight)
	proc.result <- report{source: proc.id, value: value}
}

// makeProcArray(P, input, result):
//
// Make an array of "OET processors" of size `P`.
//
// These workers engage in a protocol that mimics the behavior of
// a processor array running Odd-Even-Transposition-Sort.
//
// They communiocate via shared channels, built here, and connecting
// them as a linear array.
//
// They each obtain an initial value from the channel `input`.
// They each report their final value from the channel `result`.
//
func makeProcArray(size int,
                   input chan int, result chan report,
                   algo func(procinfo),
                   done *sync.WaitGroup) {
	var inL  chan int
    var outL chan int
	var inR  chan int = nil
	var outR chan int = nil
	for procid := 0; procid < size; procid++ {
		inL  = outR
		outL = inR
        if procid == size-1 { inR  = nil } else { inR  = make(chan int, 1)}
        if procid == size-1 { outR = nil } else { outR = make(chan int, 1)}
        info := procinfo{id:procid, procs:size,
                         inLeft:inL, outLeft:outL,
                         inRight:inR, outRight:outR,
                         input:input, result:result,
                         signal:done}
		go algo(info)
	}
}
// output(values):
//
// Output an array of values using a width of two decimal places.
func output2d(values []int) {
	for _, v := range values {
		fmt.Printf("%02d ",v)
	}
	fmt.Println()
}
