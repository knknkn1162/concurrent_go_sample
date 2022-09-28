package main

import (
    "fmt"
    "net/http"
    "golang.org/x/sync/errgroup"
)

func main() {
    urls := []string{"https://www.google.com", "https://badhost", "https://badhost2", "https://www.google.com"}
    var eg errgroup.Group
    ch := make(chan *http.Response)
    for _, url := range urls {
        url := url
        eg.Go(func() error {
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
    // onErrorEnd -> onSuccess terminate -> onErr
    if err := eg.Wait(); err != nil {
        fmt.Printf("- %v\n", err)
    }
    fmt.Println("ok")
}
