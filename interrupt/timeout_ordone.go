package main

import (
    "fmt"
    "time"
)

func orDone(done <-chan interface{}, c <-chan time.Time) <-chan interface{} {
    valStream := make(chan interface{})
    go func() {
        defer close(valStream)
        for {
            select {
            case <-done:
                return
            case v, ok := <-c:
                if ok == false {
                    return
                }
                select {
                case valStream <- v:
                case <-done:
                }
            }
        }
    }()

    return valStream
}

func do_calc() {
    for {
        fmt.Print(".")
        time.Sleep(300 * time.Millisecond)
    }
}

func main() {
    done := make(chan interface{})
    go func() {
        do_calc()
        close(done)
    }()
    // both block
    timer := time.After(3 * time.Second)
    <-orDone(done, timer)
    fmt.Println("ok")
}
