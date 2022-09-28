package main

import (
    "fmt"
    "time"
    "context"
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

func do_job_timeout(ctx context.Context) <-chan bool{
    res := make(chan bool)
    out := do_job()
    go func() {
        defer close(res)
        select {
        case <-ctx.Done():
            fmt.Println("ctx done")
        case <-out:
            fmt.Println("res")
        }
    }()
    return res
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(),3 * time.Second)
    // タイムアウト設定をしていた場合にも、明示的にcancelを呼ぶべき
    defer cancel()
    ch := do_job_timeout(ctx)
    select {
    case _, ok := <-ch:
        if ok {
            fmt.Println("return")
        } else {
            fmt.Println("timeout")
        }
    }
    fmt.Println("done!")
}
