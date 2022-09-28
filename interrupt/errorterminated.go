package main

import (
    "fmt"
    "net/http"
)

type Result struct {
    Error    error
    Response *http.Response
}

func do_job(url string) Result {
    resp, err := http.Get(url)
    return Result{Error: err, Response: resp}
}

func checkStatus(done <-chan interface{}, urls ...string) <-chan Result {
    results := make(chan Result)
    go func() {
        defer close(results)
        for _, url := range urls {
            result := do_job(url)
            select {
            case <-done:
                return
            case results <- result:
            }
        }
    }()
    return results
}

// error handler
func stopAtOnce(errCh <-chan error) <-chan bool{
    terminated := make(chan bool)
    go func() {
        defer close(terminated)
        err := <-errCh
        fmt.Printf("error occurs: %v\n", err)
        terminated<-true
    }()
    return terminated
}

func main() {
    done := make(chan interface{})
    urls := []string{"https://www.google.com", "https://badhost", "https://www.google.com"}
    errCh := make(chan error)
    defer close(errCh)
    go func() {
        defer close(errCh)
        for result := range checkStatus(done, urls...) {
            if result.Error != nil {
                errCh<-result.Error
                continue
            }
            fmt.Printf("Response: %v\n", result.Response.Status)
        }
    }()
    terminated := stopAtOnce(errCh)
    <-terminated
}
