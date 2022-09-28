package main

import (
    "fmt"
    "time"
)

func do_job(done <-chan bool, num int, c string) <-chan int{
    res := make(chan int)
    go func() {
        defer close(res)
        for i:= 0; i < num; i++ {
            select {
                case <-done:
                    return
                default:
                    fmt.Printf(c)
                    time.Sleep(300 * time.Millisecond)
                    res <-i
            }
        }
    }()
    return res
}

func main() {
    done := make(chan bool)
    terminated := make(chan bool)
    // wait
    go func() {
        terminated <-true
    }()
    res := do_job(done, 100, ".")
    for v := range res {
        if(v >= 5) {
            break
        }
    }
    // notify child to done the job!
    close(done)
    fmt.Println("press Ctrl-C")
    <-terminated
}
