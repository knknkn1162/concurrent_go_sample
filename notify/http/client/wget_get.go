package main

import (
    "net"
    "net/http"
    _ "net/http/httputil"
    "net/url"
    "fmt"
    "bufio"
)

type RequestInfo struct {
   conn net.Conn
   req *http.Request
}

func str2url(str ...string) <-chan *url.URL {
    ret := make(chan *url.URL)
    go func() {
        defer close(ret)
        for _, s := range str {
            u, err := url.Parse(s)
            if err != nil {
                fmt.Println(err)
                ret <- nil
                continue
            }
            ret <- u
        }
    }()
    return ret
}

func getConn(uCh <-chan *url.URL) <-chan net.Conn {
    ret := make(chan net.Conn)
    go func() {
        defer close(ret)
        for u := range uCh {
            conn, err := net.Dial("tcp", u.Host + u.Path + ":80")
            if err != nil {
                fmt.Println(err)
                ret <- nil
                continue
            }
            ret <- conn
        }
    }()
    return ret
}

func newreq(urlCh <-chan *url.URL) <-chan *http.Request {
    reqCh := make(chan *http.Request)
    go func() {
        defer close(reqCh)
        for u := range urlCh {
            req, err := http.NewRequest(
                "GET", u.String(), nil,
            )
            if err != nil {
                fmt.Println(err)
                reqCh <- nil
                continue
            }
            reqCh <- req
        }
    }()
    return reqCh
}

func zipRequest(connCh <-chan net.Conn, reqCh <-chan *http.Request) <-chan RequestInfo {
    ret := make(chan RequestInfo)
    go func() {
        defer close(ret)
        for req := range reqCh {
            conn := <-connCh
            if req != nil && conn != nil {
                ret <- RequestInfo{conn, req}
            }
        }
    }()
    return ret
}

func getResponse(reqInfoCh <-chan RequestInfo) <-chan *http.Response {
    resCh := make(chan *http.Response)
    go func() {
        defer close(resCh)
        for reqInfo := range reqInfoCh {
            //fmt.Printf("send request :%v\n", reqInfo.req.URL)
            reqInfo.req.Write(reqInfo.conn)
            res, err := http.ReadResponse(
                bufio.NewReader(reqInfo.conn), reqInfo.req,
            )
            if err != nil {
                fmt.Println(err)
                continue
            }
            resCh <- res
        }
    }()
    return resCh
}

func dupChannel(srcCh <-chan *url.URL) (<-chan *url.URL, <-chan *url.URL) {
    dst1 := make(chan *url.URL)
    dst2 := make(chan *url.URL)
    go func() {
        defer close(dst1)
        defer close(dst2)
        for val := range srcCh {
            dst1 <- val
            dst2 <- val
        }
    }()
    return dst1, dst2
}


func wget(urls ...string) <-chan *http.Response {
    urlCh := str2url(urls...)
    urlCh1, urlCh2 := dupChannel(urlCh)
    connCh := getConn(urlCh1)
    reqCh := newreq(urlCh2)
    reqInfoCh := zipRequest(connCh, reqCh)
    return getResponse(reqInfoCh)
}

func printResponse(resCh <-chan *http.Response) <-chan bool {
    ret := make(chan bool)
    go func() {
        defer close(ret)
        for res := range resCh {
            fmt.Printf("%v %v\n", res.Proto, res.Status)
        }
    }()
    return ret
}

func main() {
    urls := []string{"https://www.google.com", "https://badhost", "https://www.google.com", "https://www.google.com"}
    resCh := wget(urls...)
    done := printResponse(resCh)
    <-done
}
