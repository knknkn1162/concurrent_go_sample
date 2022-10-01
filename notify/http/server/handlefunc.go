package main

import (
    "net/http"
    "io"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "sample")
    })
    http.ListenAndServe(":8080", nil)
}
