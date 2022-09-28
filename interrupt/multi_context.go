package main

import (
    "fmt"
    "time"
    "context"
    "errors"
    "math/rand"
)

func do_calc() int {
    for {
        time.Sleep(300 * time.Millisecond)
        fmt.Print(".")
    }
    return 1
}

func do_job() <-chan int {
    out := make(chan int)
    go func(){
        defer close(out)
        out<-do_calc()
    }()
    return out
}

func do_job_with(ctx context.Context) <-chan bool{
    res := make(chan bool)
    out := do_job()
    go func() {
        defer close(res)
        select {
        case <-ctx.Done():
            if err := ctx.Err(); errors.Is(err, context.Canceled) {
                fmt.Println("canceled")
            } else if errors.Is(err, context.DeadlineExceeded) {
                fmt.Println("deadline")
            }
        case <-out:
            fmt.Println("res")
        }
    }()
    return res
}

func main() {
    ctx1, cancelTimeout := context.WithTimeout(context.Background(),3 * time.Second)
    ctx, cancel := context.WithCancel(ctx1)
    // タイムアウト設定をしていた場合にも、明示的にcancelを呼ぶべき
    defer cancelTimeout()
    defer cancel()
    ch := do_job_with(ctx)
    // 0: canceled 2: deadline
    rand.Seed(2)
LOOP:
    for {
        select {
        case _, ok := <-ch:
            if ok {
                fmt.Println("return")
            } else {
                fmt.Println("stopped")
            }
            break LOOP
        default:
            if(rand.Intn(1000000004) >= 1000000000) {
                cancel()
            }
        }
    }

    fmt.Println("done!")
}
