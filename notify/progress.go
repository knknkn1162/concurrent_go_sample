package main

import (
    "fmt"
    "time"
    "math/rand"
)

func check_progress(done1 <-chan bool) {
    for {
        <-done1
        fmt.Printf(".")
    }
}

func do_calc(done1 chan<- bool) {
    for {
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
        done1<-true
    }
}

func main() {
    done1 := make(chan bool)
    go func() {
        do_calc(done1)
    }()
    check_progress(done1)
}
