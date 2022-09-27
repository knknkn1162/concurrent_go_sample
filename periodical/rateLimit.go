package main

import (
    "fmt"
    "time"
    "math/rand"
)


func do_calc(id int) {
    n := rand.Intn(4000)
    time.Sleep(time.Duration(n) * time.Millisecond)
    fmt.Printf("end %v(%v)\n", id, n)
}

func main() {
    ticker := time.Tick(time.Second)
    for i:= 0; ; i++{
        <-ticker
        fmt.Printf("start %v\n", i)
        go do_calc(i)
    }
}
