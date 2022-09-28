package main

import (
    "fmt"
    "time"
)

func do_calc(num int, c string) {
    for i:= 0; i < num; i++ {
        fmt.Printf(c)
        time.Sleep(300 * time.Millisecond)
    }
}

func printFinished(done <-chan bool) <-chan bool{
    terminated := make(chan bool)
    go func() {
        defer close(terminated)
        for {
            select {
            case <-done:
                fmt.Println("finished!")
                return
            }
        }
    }()
    return terminated
}

func main() {
    done := make(chan bool)
    terminated := printFinished(done)
    // finished
    go func() {
        do_calc(5, ".")
        close(done)
    }()
    <-terminated
}
