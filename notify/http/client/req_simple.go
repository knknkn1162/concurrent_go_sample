package main

import (
    "net/http"
    "fmt"
)

func check_req(urls []string) <-chan *http.Request {
    ch := make(chan *http.Request)
    go func() {
        defer close(ch)
        for _, url := range urls {
            req, err := http.NewRequest(
                "GET", url, nil,
            )
            if err != nil {
                return
            }
            ch <- req
        }
    }()
    return ch
}

func main() {
    urls := []string{"https://www.google.com", "https://badhost", "https://www.google.com", "https://www.google.com"}
    for req := range check_req(urls) {
        fmt.Printf("url: %v (%v)\n", req.URL, req.Proto)
    }
}
