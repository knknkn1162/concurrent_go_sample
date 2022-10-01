package main

import (
    "os"
    "io"
    "bytes"
    "fmt"
    "encoding/binary"
)

func main() {
    fpath := "./file.txt"
    file, err := os.Open(fpath)
    if err != nil {
        panic(err)
    }
    buf, err := io.ReadAll(file)
    if err != nil {
        panic(err)
    }

    byteReader := bytes.NewReader(buf)
    byteLen := byteReader.Len()
    fmt.Printf("total %v\n", byteLen)
    chunks := make(chan []byte)

    go func() {
        defer close(chunks)
        for {
            unit := 23
            if (byteLen < unit) {
                unit = byteLen
            }
            chunk := make([]byte, unit)
            err := binary.Read(byteReader, binary.BigEndian, &chunk)

            chunks <- chunk
            byteLen -= unit
            if byteLen == 0 {
                break
            }
            if err == io.EOF {
                break
            }
        }
    }()

    done := make(chan bool)
    go func() {
        defer close(done)
        for bytes := range chunks {
            fmt.Printf("%v (len: %v)\n", bytes, len(bytes))
        }
    }()
    <-done
    fmt.Println("finished!")
}
