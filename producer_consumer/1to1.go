package main

import (
    "fmt"
    "time"
)

func consume(jobs <-chan int) <-chan int{
    done := make(chan int)
    go func() {
        defer close(done)
        for j := range jobs {
            fmt.Printf("started  job %v", j)
            time.Sleep(time.Second)
            fmt.Println(" => finished job")
        }
    }()
    return done
}

func produce() <-chan int{
    jobs := make(chan int)
    go func() {
        defer close(jobs)
        for j := 1; j <= 5; j++ {
            jobs <- j
        }
    }()
    return jobs
}

func main() {
    // 1
    jobs := produce()
    <-consume(jobs)
}
