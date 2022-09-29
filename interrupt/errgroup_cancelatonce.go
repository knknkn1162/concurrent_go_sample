package main

import (
    "fmt"
    "net/http"
    "golang.org/x/sync/errgroup"
    "context"
    "time"
)

func main() {
    urls := []string{"https://www.google.com", "https://badhost", "https://badhost2", "https://www.google.com", "https://www.google.com", "https://www.google.com", "https://www.google.com", "https://www.google.com"}

    eg, ctx := errgroup.WithContext(context.Background())
    ch := make(chan *http.Response)
    for _, url := range urls {
        url := url
        // little harder
        time.Sleep(20 * time.Millisecond)
        eg.Go(func() error {
            // early return
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }
            resp, err := http.Get(url)
            if err != nil {
                return err
            }
            ch <- resp
            return nil
        })
    }

    // onErrorEnd
    go func() {
        err := eg.Wait()
        fmt.Printf("onEnd: %v\n", err)
        close(ch)
    }()

    // onSuccess
    for response := range ch {
        fmt.Printf("+ %v\n", response.Status)
    }
    // onFailure
    if err := eg.Wait(); err != nil {
        fmt.Printf("- %v\n", err)
    }
    fmt.Println("ok")
}
