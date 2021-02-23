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
    round int
	sourcei int
    sourcej int
	value  int
}

// procinfo struct
//
// Info pertaining to a "processor"'s configuration within its "linear array".
type meshinfo struct {
    row, rows int
    column, columns int
    inNorth, outNorth chan int
    inEast, outEast chan int
    inSouth, outSouth chan int
    inWest, outWest chan int
    input chan int
    result chan report
    signal *sync.WaitGroup
}

func lg(n int) int {
    i := 0
    for n > 1 {
        n /= 2
        i++
    }
    return i
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
	rows := 4
    cols := 4
	if len(os.Args) > 2 {
		rs, err := strconv.Atoi(os.Args[1])
		if err == nil {
			rows = rs
		}
        cs, err := strconv.Atoi(os.Args[2])
        if err == nil {
            cols = cs
        }
	}
    size := rows*cols

	// Create a random permutation of values.
	rand.Seed(time.Now().UnixNano())
	values := randomPerm(size)

    steps := lg(rows)*2 + 1
    fmt.Println(steps)
    log := make([][][]int,steps)
    for t := 0; t < steps; t++ {
        log[t] = make([][]int,rows)
        for i := 0; i < rows; i++ {
            log[t][i] = make([]int,cols)
            for j := 0; j < cols; j++ {
                log[t][i][j] = 99
            }
        }
    }


	// Create channels to ship values to/from the proc array.
	inputs := make(chan int, size)
	results := make(chan report, size*steps)

	// Print unsorted data.
	output3d(values)

    var done sync.WaitGroup
    done.Add(size);

	// Create the processor array to perform oetSort.
	makeProcMesh(rows, cols, inputs, results, shearSort, &done)

	// Scatter the data to them.
	for v := range values {
		inputs <- v
	}

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
	    // Gather the sorted data from them.
	    for count := 0; count < steps*rows*cols; count++ {
            sd := <-results
		    log[sd.round][sd.sourcei][sd.sourcej] = sd.value
	    }
    }()

    done.Wait()
    wg.Wait()

	// Print sorted data.
	output3d(values)

    for s := 0; s < steps; s++ {
        fmt.Println("-------")
        fmt.Println("Step ",s)
        for i := 0; i < rows; i++ {
            output3d(log[s][i])
        }
    }
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

func snake(value *int, step int, proc meshinfo) {
    if proc.row % 2 == 0 {
        oetRun(value, proc.column, proc.columns,
               proc.outWest, proc.inWest,
               proc.outEast, proc.inEast)
        proc.result <- report{round:step,
                              sourcei:proc.row, sourcej:proc.column,
                              value:*value}
    } else {
        oetRun(value, proc.columns-proc.column-1, proc.columns,
               proc.outEast, proc.inEast,
               proc.outWest, proc.inWest)
        proc.result <- report{round:step,
                              sourcei:proc.row, sourcej:proc.column,
                              value:*value}
    }
}

func up(value *int, step int, proc meshinfo) {
    oetRun(value, proc.row, proc.rows,
           proc.outNorth, proc.inNorth,
           proc.outSouth, proc.inSouth)
    proc.result <- report{round:step,
                          sourcei:proc.row, sourcej:proc.column,
                          value:*value}
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
func shearSort(proc meshinfo) {
    defer proc.signal.Done()
	value := <-proc.input

    t := 0
    for phase := proc.rows; phase > 1; phase /= 2 {
        snake(&value,t,proc)
        up(&value,t+1,proc)
        t += 2
    }
    snake(&value,t,proc)

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
func makeProcMesh(rows int, cols int,
                  input chan int, result chan report,
                  algo func(meshinfo),
                  done *sync.WaitGroup) {
    var inN  []chan int
    var outN []chan int
    var inS  []chan int = make([]chan int,cols)
    var outS []chan int = make([]chan int,cols)
    for proci := 0; proci < rows; proci++ {

        inN  = outS
        outN = inS
        inS  = make([]chan int,cols)
        outS = make([]chan int,cols)
        if proci != rows-1 {
            for j := 0; j < cols; j++ {
                inS[j]  = make(chan int,1)
                outS[j] = make(chan int,1)
            }
        }

    	var inW  chan int
        var outW chan int
    	var inE  chan int = nil
    	var outE chan int = nil

    	for procj := 0; procj < cols; procj++ {
    		inW  = outE
    		outW = inE
            if procj == cols-1 { inE  = nil } else { inE  = make(chan int, 1)}
            if procj == cols-1 { outE = nil } else { outE = make(chan int, 1)}
            info := meshinfo{row:proci, rows:rows,
                             column:procj, columns:cols,
                             inNorth:inN[procj], outNorth:outN[procj],
                             inEast:inE, outEast:outE,
                             inSouth:inS[procj], outSouth:outS[procj],
                             inWest:inW, outWest:outW,
                             input:input, result:result,
                             signal:done}
    		go algo(info)
    	}
    }
}
// output(values):
//
// Output an array of values using a width of two decimal places.
func output3d(values []int) {
	for _, v := range values {
		fmt.Printf("%03d ",v)
	}
	fmt.Println()
}
