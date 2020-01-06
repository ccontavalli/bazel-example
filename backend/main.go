package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/ccontavalli/bazel-example/backend/lib"
)

func main() {
    http.HandleFunc("/", HelloServer)
    log.Printf("Opening port 5432 - will be available at http://127.0.0.1:5432/")
    err := http.ListenAndServe(":5432", nil)
    log.Printf("Listen returned error - %v", err)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    lib.MyLog(r.URL.Path)
}
