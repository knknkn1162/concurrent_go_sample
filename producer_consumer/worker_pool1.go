package main

import (
    "fmt"
    "time"
)

func consume(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("consumer", id, "started  job", j)
        time.Sleep(time.Second)
        fmt.Println("consumer", id, "finished job", j)
        results <- j * 2
    }
}

func produce(jobs chan int) {
    defer close(jobs)
    for j := 1; j <= 10; j++ {
        jobs <- j
    }
}
func main() {
    const numJobs = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    // 1
    go produce(jobs)

    // n: harder
    for w := 1; w <= 3; w++ {
        go consume(w, jobs, results)
    }

    for a := 1; a <= numJobs; a++ {
        res := <-results
        fmt.Printf("receive %v\n", res)
    }
}
