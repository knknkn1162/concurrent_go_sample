package main

import (
  "fmt"
  "time"
)

func worker(id int) {
  fmt.Printf("Worker %d starting\n", id)
  defer fmt.Printf("Worker %d done\n", id)
  time.Sleep(time.Second)
}

func finished(num int, done <-chan int) <-chan bool{
    terminated := make(chan bool)
    var cnt uint64
    go func() {
        defer close(terminated)
        for {
            select {
                case id := <-done:
                    fmt.Printf("job: %v done\n", id)
                    cnt++
            }
            if(int(cnt) >= num) {
                break
            }
        }
    }()
    return terminated
}

func main() {
    done := make(chan int)
    const num = 5
    for i := 1; i <= num; i++ {
      go func(i int) {
        worker(i)
        done<-i
      }(i)
    }
    terminated := finished(num, done)
    <-terminated
    fmt.Println("finished!")
}
