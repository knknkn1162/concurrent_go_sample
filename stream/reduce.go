package main

import (
    "fmt"
)

func gen(integers ...int) <-chan int {
    istream := make(chan int)
    go func() {
        defer close(istream)
        for _, i := range(integers) {
            istream <- i
        }
    }()
    return istream
}

func reduce(stream <-chan int, additive int) <-chan int {
    res := make(chan int)
    var ans int
    go func() {
        defer close(res)
        for v := range stream {
            ans += v
        }
        res <- ans
    }()
    return res
}

func main() {
    ans := reduce(gen(1,2,3,4), 1)
    fmt.Printf("ans: %v\n", <-ans)
}
