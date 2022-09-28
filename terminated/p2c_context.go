// parent -> child nnotification
package main

import (
    "fmt"
    "time"
    "context"
)

func do_calc(num int, c string) int {
    fmt.Printf(c)
    time.Sleep(300 * time.Millisecond)
    return 1
}

func do_job(ctx context.Context, num int, c string) <-chan int{
    res := make(chan int)
    go func() {
        defer close(res)
        for i:= 0; i < num; i++ {
            select {
                case <-ctx.Done():
                    return
                default:
                    res <-do_calc(i, c)
            }
        }
    }()
    return res
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    res := do_job(ctx, 100, ".")
    cnt := 0
    for v := range res {
        cnt += v
        if(cnt >= 5) {
            break
        }
    }
    // notify child to done the job!
    cancel()
    fmt.Println("done!")
}
