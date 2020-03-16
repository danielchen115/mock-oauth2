package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", ServerHandler)
    http.ListenAndServe(":7080", nil)
}

func ServerHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Mock OAuth 2 Server v1")
}
