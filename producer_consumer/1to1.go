package main

import (
    "fmt"
    "time"
)

func consume(jobs <-chan int, results chan<- int) {
    defer close(results)
    for j := range jobs {
        fmt.Println("consumer", "started  job")
        time.Sleep(time.Second)
        fmt.Println("consumer", "finished job")
        results <- j * 2
    }
}

func produce(jobs chan int) {
    defer close(jobs)
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
}
func main() {
    jobs := make(chan int)
    results := make(chan int)
    // 1
    go produce(jobs)
    // 1
    go consume(jobs, results)

    for n := range results {
        fmt.Printf("receive %v\n", n)
    }
}
