package main

import (
    "fmt"
    "time"
)

func do_calc() {
    fmt.Print(".")
    time.Sleep(300 * time.Millisecond)
}

func main() {
    timeout := time.After(3 * time.Second)
    for {
        select {
            case <-timeout:
                fmt.Println("timeout")
                return
            default:
                do_calc()
        }
    }
}
