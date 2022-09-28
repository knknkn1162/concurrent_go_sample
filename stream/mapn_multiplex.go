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

func add(stream <-chan int, additive int) []chan int {
    const maxCh = 3
    ans := make([]chan int, 3)
    for i := 0; i < maxCh; i++ {
        res := make(chan int)
        ans[i] = res
        go func() {
            defer close(res)
            for v := range stream {
                // it's little heavy
                res <- do_calc(v, additive)
            }
        }()
    }
    return ans
}

func fanin(chs []chan int) <-chan int{
    res := make(chan int)
    var wg sync.WaitGroup
    // onClose
    go func() {
        defer close(res)
        wg.Wait()
    }()

    for _, ch := range chs {
        wg.Add(1)
        go func(ch chan int) {
            defer wg.Done()
            for i := range ch {
                res <- i
            }
        }(ch)
    }
    return res
}

func main() {
    // []chan int
    chs := add(gen(15), 1)
    res := fanin(chs)
    for v := range res {
        fmt.Printf("ans: %v\n", v)
    }
}
