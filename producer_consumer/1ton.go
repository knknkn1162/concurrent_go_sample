package main

import (
    "fmt"
    "time"
    "sync"
)

const numJobs = 5
func consume(jobs <-chan int) <-chan int {
    results := make(chan int, numJobs)
    var wg sync.WaitGroup
    for w := 1; w <= 3; w++ {
        wg.Add(1)
        go func(w int) {
            defer wg.Done()
            // it takes a little time
            for j := range jobs {
                fmt.Println("consumer", w, "started  job", j)
                time.Sleep(time.Second)
                fmt.Println("consumer", w, "finished job", j)
                results <- j * 2
            }
        }(w)
    }
    // onEnd
    go func() {
        defer close(results)
        wg.Wait()
    }()
    return results
}

func produce() <-chan int {
    jobs := make(chan int, numJobs)
    go func() {
        defer close(jobs)
        for j := 1; j <= 10; j++ {
            fmt.Printf("generate %v\n", j)
            jobs <- j
        }
    }()
    return jobs
}

func main() {
    // 1
    jobCh := produce()
    resultCh := consume(jobCh)
    for res := range resultCh {
        fmt.Printf("receive %v\n", res)
    }
}
