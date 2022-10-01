package main

import (
    "fmt"
    "net"
    "sync"
    "net/url"
    "net/http"
    "bufio"
    "strings"
    "io"
    "strconv"
    "crypto/sha256"
)

func accept(jobNum int, listener net.Listener) <-chan net.Conn {
    ret := make(chan net.Conn)
    var wg sync.WaitGroup
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for {
                fmt.Printf("accept @%v\n", i)
                conn, err := listener.Accept()
                if err != nil {
                    panic(err)
                }
                ret <- conn
            }
        }(i)
    }
    go func() {
        defer close(ret)
        wg.Wait()
    }()
    return ret
}

func url2data(url *url.URL) int {
    qs := url.Query()
    // extract data
    num, err := strconv.Atoi(qs.Get("repeat"))
    if err != nil {
        num = 0
    }
    return num
}

// do little heavy job
// note that you should do the other computer
func do_calc(num int) string {
    intCh := make(chan int)
    go func() {
        defer close(intCh)
        for i := 0; i < num; i++ {
            intCh <- i
        }
    }()

    strCh := make(chan string)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            for val := range intCh {
                strCh <- string(val)
            }
        }(i)
    }
    go func() {
        defer close(strCh)
        wg.Wait()
    }()
    h := sha256.New()
    for val := range strCh {
        h.Write([]byte(val))
    }
    sha256 := h.Sum(nil)
    return fmt.Sprintf("%v: 0x%x", num, sha256) + "\n"
}

func sendData(conn net.Conn, idx int) {
    defer conn.Close()
    fmt.Printf("recv %v\n", idx)
    req, err := http.ReadRequest(
        bufio.NewReader(conn),
    )
    if err != nil {
        panic(err)
    }
    num := url2data(req.URL)
    content := do_calc(num)

    response := http.Response{
        StatusCode: 200,
        Proto: "1.1",
        ContentLength: int64(len(content)),
        Body: io.NopCloser(strings.NewReader(content)),
    }
    response.Write(conn)
}

func readRequest(jobNum int, connCh <-chan net.Conn) <-chan bool {
    done := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            for conn := range connCh {
                sendData(conn, idx)
            }
        }(i)
    }
    go func() {
        defer close(done)
        wg.Wait()
    }()
    return done
}

func main() {
    listener, err := net.Listen("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }
    connCh := accept(5, listener)
    <-readRequest(30, connCh)
    fmt.Println("ok")
}
