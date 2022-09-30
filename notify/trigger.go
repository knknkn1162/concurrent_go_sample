package main

import (
    "os"
    "os/signal"
    "fmt"
    "time"
)

func do_calc(i int) int {
    fmt.Printf("job #%v", i)
    time.Sleep(300 * time.Millisecond)
    fmt.Printf(" => ok\n")
    return i * 2
}

// 単一の場合
func do_job(i int) <-chan int{
    res := make(chan int)
    go func() {
        defer close(res)
        time.Sleep(100 * time.Millisecond)
        res<-do_calc(i)
    }()
    return res
}

func usesig() <-chan bool{
    done := make(chan bool)
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    go func() {
        defer close(done)
        defer close(sigint)
        cnt := 0
        for {
            // trigger
            <-sigint
            <-do_job(cnt)
            cnt++;
        }
    }()
    return done
}

func main() {
    done := usesig()
    <-done
}
