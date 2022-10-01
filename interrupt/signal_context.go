package main
import (
    "context"
    "fmt"
    "os/signal"
    "syscall"
    "time"
)

func do_calc() {
    for {
        fmt.Print(".")
        time.Sleep(1 * time.Second)
    }
}

func main() {
    ctx := context.Background()
    sigctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    go func() {
        do_calc()
    }()
    <-sigctx.Done()
    fmt.Println("shutdown")
}
