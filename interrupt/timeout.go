package main

import (
    "fmt"
    "time"
)

func do_calc() {
    for {
        fmt.Print(".")
        time.Sleep(300 * time.Millisecond)
    }
}

func main() {
    done := make(chan bool)
    go func() {
        do_calc()
        done<-true
    }()
    // both block
    select {
        case <-done:
        case <-time.After(3 * time.Second):
            fmt.Println("timeout")
    }
}
