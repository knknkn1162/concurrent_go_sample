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

func add(stream <-chan int, additive int) <-chan int {
    res := make(chan int)
    go func() {
        defer close(res)
        for v := range stream {
            res <- v + additive
        }
    }()
    return res
}

func main() {
    pp := add(gen(1,2,3,4), 1)
    for v := range pp {
        fmt.Printf("ans: %v\n", v)
    }
}
