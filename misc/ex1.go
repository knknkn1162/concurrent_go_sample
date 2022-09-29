package main

import (
    "context"
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "golang.org/x/sync/errgroup"
)

func main() {
    ctx := context.Background()
    pathStream := gen_path(ctx, "..")
    //pathStream := gen_path_fake()
    resultStream := file2md5(ctx, pathStream)
    for res := range resultStream {
        fmt.Printf("%s:\t%x\n", res.path, res.sum)
    }
}

type result struct {
    path string
    sum  [md5.Size]byte
}

func gen_path_fake() <-chan string {
    paths := make(chan string)
    go func() {
        defer close(paths)
        paths <- "ssm"
    }()
    return paths

}

func gen_path(ctx context.Context, root string) <-chan string {
    paths := make(chan string)
    go func() {
        defer close(paths)
        filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.Mode().IsRegular() {
                return nil
            }
            select {
            case paths <- path:
            case <-ctx.Done():
                return ctx.Err()
            }
            return nil
        })
    }()
    return paths
}

func file2md5(ctx context.Context, paths <-chan string) <-chan result {
    g, ctx := errgroup.WithContext(ctx)
    ch := make(chan result)
    const numDigesters = 20
    // parallel
    for i := 0; i < numDigesters; i++ {
        g.Go(func() error {
            for path := range paths {
                data, err := ioutil.ReadFile(path)
                if err != nil {
                    return err
                }
                select {
                case ch <- result{path, md5.Sum(data)}:
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
            return nil
        })
    }
    // onError
    go func() {
        defer close(ch)
        err := g.Wait()
        if err != nil {
            fmt.Printf("err: %v\n", err)
        } else {
            fmt.Println("completed!")
        }
    }()
    return ch
}
