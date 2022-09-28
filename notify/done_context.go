package main
import (
    "sync"
    "fmt"
    "context"
)
var wg sync.WaitGroup


func main() {
    ctx, cancel := context.WithCancel(context.Background())
    gen := generator(ctx, 1)

    wg.Add(1)

    for i := 0; i < 5; i++ {
        fmt.Println(<-gen)
    }
    cancel()

    wg.Wait()
}
