//go:build profile

package main

import (
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()
}
