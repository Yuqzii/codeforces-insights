package main

import (
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/transport"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", transport.HandleRoot)

	fmt.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", mux)
}
