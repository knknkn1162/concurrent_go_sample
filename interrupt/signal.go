package main

import (
    "os"
    "os/signal"
    "fmt"
    "time"
)

func do_calc() {
    for {
        fmt.Print(".")
        time.Sleep(1 * time.Second)
    }
}

func watchSignal(cancel chan bool) {
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt)
    <-sigint
    cancel<-true
}

func main() {
    cancel := make(chan bool)
    go watchSignal(cancel)
    go func() {
        do_calc()
    }()
    // wait until cancel=true
    <-cancel
    fmt.Println("canceled!")
}
