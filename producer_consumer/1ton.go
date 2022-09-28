package main

import (
    "fmt"
    "time"
)

const numJobs = 5
func consume(id int, jobs <-chan int, results chan<- int) {
    defer close(results)
    for j := range jobs {
        fmt.Println("consumer", id, "started  job", j)
        time.Sleep(time.Second)
        fmt.Println("consumer", id, "finished job", j)
        results <- j * 2
    }
}

func produce() chan int {
    jobs := make(chan int, numJobs)
    go func() {
        defer close(jobs)
        for j := 1; j <= 10; j++ {
            jobs <- j
        }
    }()
    return jobs
}
func main() {
    // 1
    jobs := produce()

    results := make(chan int, numJobs)
    // n: harder
    for w := 1; w <= 3; w++ {
        go consume(w, jobs, results)
    }

    for res := range results {
        fmt.Printf("receive %v\n", res)
    }
}
