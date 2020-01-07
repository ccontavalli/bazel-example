package main

import (
    "fmt"
    "net/http"
    "log"

    // Includes an example library.
    "github.com/ccontavalli/bazel-example/backend/lib"
    // Includes an example set of assets, created with go_embed_data.
    "github.com/ccontavalli/bazel-example/backend/assets"
)

func main() {
    http.HandleFunc("/", HelloServer)
    log.Printf("Opening port 5432 - will be available at http://127.0.0.1:5432/")
    log.Printf("Serving assets:")
    for key, data := range assets.Data {
      log.Printf("  - %s - %d bytes", key, len(data))
    }
    err := http.ListenAndServe(":5432", nil)
    log.Printf("Listen returned error - %v", err)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    lib.MyLog(r.URL.Path)
}
