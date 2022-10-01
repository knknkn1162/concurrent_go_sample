package main

import (
    "io"
    "fmt"
    "net"
    "net/http"
    "sync"
    "strings"
    "bufio"
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

func handleRequest(req *http.Request) <-chan *http.Response {
    sessionResponse := make(chan *http.Response)
    go func() {
        defer close(sessionResponse)
        content := "hello world\n"
        resp := &http.Response{
            StatusCode: 200,
            Proto: "1.1",
            ContentLength: int64(len(content)),
            Body: io.NopCloser(strings.NewReader(content)),
        }
        // 2: send data in sessionResponse
        sessionResponse <- resp
        // 4: close sessionResponse
    }()
    return sessionResponse
}

func sendData(conn net.Conn, idx int) <-chan<-chan *http.Response {
    fmt.Printf("recv %v\n", idx)
    sessionResponses := make(chan(<-chan *http.Response))
    go func() {
        defer close(sessionResponses)
        reader := bufio.NewReader(conn)
        for {
            req, err := http.ReadRequest(reader)
            if err != nil {
                //neterr, ok := err.(net.Error)
                if err == io.EOF {
                    fmt.Println("io.EOF")
                    break
                }
                panic(err)
            }
            fmt.Println("publish sessionResponse")
            sessionResponse := handleRequest(req)
            // 1: send sessionResponse
            sessionResponses<- sessionResponse
        }
        // 5: close sessionResponses
    }()
    return sessionResponses
}
func writeToConn(sessionResponses <-chan<-chan *http.Response, conn net.Conn) chan interface{}{
    done := make(chan interface{})
    go func() {
        defer close(done)
        defer conn.Close()
        for sessionResponse := range sessionResponses {
            fmt.Println("recv: sessionResponse")
            // 3: recv data in sessionResponse
            // // 値が残っていたらそれを、なければデフォルト値を返すので、ここはnilにならない
            resp := <-sessionResponse
            resp.Write(conn)
        }
        // 6: end loop for sessionResponses
        fmt.Println("end: sessionResponses")
    }()
    return done
}

func readRequest(jobNum int, connCh <-chan net.Conn) <-chan bool {
    done := make(chan bool)
    var wg sync.WaitGroup
    for i := 0; i < jobNum; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            for conn := range connCh {
                sessionResponses := sendData(conn, idx)
                // async
                writeToConn(sessionResponses, conn)
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
