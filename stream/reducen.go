package main

import (
    "fmt"
    "sync/atomic"
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

func reduce(stream <-chan int, additive int) <-chan int {
    res := make(chan int)
    var ans int32
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for v := range stream {
                // must use atomic function
                atomic.AddInt32(&ans, int32(v))
            }
        }()
    }
    go func() {
        wg.Wait()
        res <- int(ans)
    }()
    return res
}

func main() {
    ans := reduce(gen(15), 1)
    fmt.Printf("ans: %v\n", <-ans)
}
