// parent -> child nnotification
package main

import (
    "fmt"
    "time"
)

func do_calc(num int, c string) int {
    fmt.Printf(c)
    time.Sleep(300 * time.Millisecond)
    return 1
}

func do_job(done <-chan bool, num int, c string) <-chan int{
    res := make(chan int)
    go func() {
        defer close(res)
        for i:= 0; i < num; i++ {
            select {
                case <-done:
                    return
                default:
                    res <-do_calc(i, c)
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
    cnt := 0
    for v := range res {
        cnt += v
        if(cnt >= 5) {
            break
        }
    }
    // notify child to done the job!
    close(done)
    fmt.Println("done!")
    <-terminated
    fmt.Println("terminated")
}
