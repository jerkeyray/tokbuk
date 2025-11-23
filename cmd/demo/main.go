package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jerkeyray/tokbuk/internal/limiter"
	"github.com/jerkeyray/tokbuk/pkg/middleware"
)

// KeyFunc implementation that identifies clients by IP address.
func clientKey(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func main() {
	manager := limiter.NewBucketManager()

	// build rate limit middleware with capacity = 10, rate = 5
	rl := middleware.RateLimit(manager, 10, 5, clientKey)

	// create router
	mux := http.NewServeMux()

	// wrap "/test" handler with rate limit middleware
	mux.Handle("/test", rl(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "allowed")
	})))

	// configure and create HTTP server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Println("middleware demo running at http://localhost:8080/test")

	// start server and exit if it fails
	log.Fatal(server.ListenAndServe())
}
