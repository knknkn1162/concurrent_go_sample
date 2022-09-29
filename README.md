# go concurrency patterns

# url

## official
+ [RateLimiting](https://github.com/golang/go/wiki/RateLimiting)
  + time.Tick(time.Second)
+ [SignalHandling](https://github.com/golang/go/wiki/SignalHandling)
  + sig := make(chan os.Signal, 1); signal.Notify(sig, os.Interrupt)
+ [Timeouts](https://github.com/golang/go/wiki/Timeouts)
  + `time.After(timeoutNanoseconds)(chan time.Time)`

## repo samples
+ https://github.com/lotusirous/go-concurrency-patterns
+ https://github.com/luk4z7/go-concurrency-guide
+ https://github.com/gohandson/goroutine-ja
