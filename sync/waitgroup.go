package main

import (
  "fmt"
  "sync"
  "time"
)

func worker(id int) {
  fmt.Printf("Worker %d starting\n", id)
  defer fmt.Printf("Worker %d done\n", id)
  time.Sleep(time.Second)
}

func main() {
  var wg sync.WaitGroup
  for i := 1; i <= 5; i++ {
    wg.Add(1)
    go func(i int) {
      defer wg.Done()
      worker(i)
    }(i)
  }
  // wait until all worker ends
  wg.Wait()
}
