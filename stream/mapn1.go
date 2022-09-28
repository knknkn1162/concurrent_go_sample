package main

import (
    "fmt"
    "time"
    "sync"
)

func gen(max int) <-chan int {
    istream := make(chan int)
    go func() {
        defer close(istream)
        for i := 0; i < max; i++ {
            istream <- i
        }
    }()
    return istream
}

func do_calc(num1 int, num2 int) int {
    time.Sleep(100 * time.Millisecond)
    return num1 + num2
}

func add(stream <-chan int, additive int) <-chan int {
    res := make(chan int)
    // when all co-goroutines are finished
    // finish this with close(res)
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for v := range stream {
                // it's little heavy
                res <- do_calc(v, additive)
            }
        }()
    }
    // when all jobs are done, close channel
    go func() {
        defer close(res)
        wg.Wait()
    }()
    return res
}

func main() {
    pp := add(gen(15), 1)
    for v := range pp {
        fmt.Printf("ans: %v\n", v)
    }
}
