package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
    "runtime/pprof"
    "runtime/trace"
    "flag"
)

func fillUp(max int, xs []int, l int, r int, wait *sync.WaitGroup) {
	defer wait.Done()
	for i := l; i < r; i++ {
		xs[i] = rand.Intn(max)
	}
}

func permute(xs []int) {
	// Generate a random permutation.
	size := len(xs)
	for i := 0; i < size; i += 1 {
		xs[i] = i
	}
	for i := 0; i < size; i += 1 {
		j := rand.Intn(size-i) + i
		xs[i], xs[j] = xs[j], xs[i]
	}
}

func sumUp(I int, ss[] int, xs []int, l int, r int, wait *sync.WaitGroup) {
	defer wait.Done()
	s := 0
	for i := l; i < r; i++ {
		s += xs[i]
	}
	ss[I] = s
}

func main() {

    // Get parameters.
    var fP = flag.Int("P", 4, "number of processors")
    var fN = flag.Int("N", 16, "amount of data to sum")
    var fperm = flag.Bool("perm", false, "use a random permutation")
    var fcpu = flag.Bool("cpu", false, "generate cpu pprof data")
    var fexe = flag.Bool("exe", false, "generate execution trace data")
    //
    flag.Parse()
    //
    P := *fP
    N := *fN
    perm := *fperm
    cpu := *fcpu
    exe := *fexe

    // Maybe generate profiling or execution traces.
    var tracename = ""
    //
    if cpu && !exe {
        tod := time.Now().Format("060102150405_pprof.out")
        tracename = "traces/"+os.Args[0]+"_"+strconv.Itoa(N)+"_"+strconv.Itoa(P)+"_"+tod
        f,_ := os.Create(tracename)
        defer func() {
            f.Close()
            fmt.Println("Wrote '"+tracename+"' file.")
        }()
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    //
    if exe {
        tod := time.Now().Format("060102150405_trace.out")
        tracename = "traces/"+os.Args[0]+"_"+strconv.Itoa(N)+"_"+strconv.Itoa(P)+"_"+tod
        f,_ := os.Create(tracename)
        defer func() {
            f.Close()
            fmt.Println("Wrote '"+tracename+"' file.")
        }()
        trace.Start(f)
        defer trace.Stop()
    }

    // Make a wait group to coordinate with any workers.
    var wait sync.WaitGroup

	// Initialize the data (using several threads).
    var values []int = make([]int, N)
    n := N / P
    rand.Seed(time.Now().UnixNano())
    if perm {
        permute(values)
    } else {
        wait.Add(P)
	    for I := 0; I < P; I++ {
		     go fillUp(100, values, I * n , I * n + n, &wait)
	    }
        wait.Wait()
    }
    if N <= 24 {
        fmt.Println("Data:", values)
    }

	// Sum up the data in parallel, timing the parallel work.
    results := make([]int, P)
	wait.Add(P)
	start := time.Now()
	for I := 0; I < P; I++ {
		go sumUp(I, results, values, I * n , I * n + n, &wait)
	}
	wait.Wait()

    // Report what happened.
	duration := time.Since(start)
  	sum := 0
    for I := 0; I < P; I++ {
   	    sum += results[I]
    }
   	fmt.Println("The sum computed was", sum)
    if perm {
        fmt.Println("The sum expected was", N * (N - 1) / 2)
    }
	fmt.Println("The time to compute it was", duration)
	fmt.Println(P, "workers ran on as mamy as", runtime.GOMAXPROCS(0), "processors.")
}
