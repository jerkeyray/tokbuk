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

func clientKey(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func main() {
	manager := limiter.NewBucketManager()

		rl := middleware.RateLimit(manager, 10, 5, clientKey)

		mux := http.NewServeMux()

		mux.Handle("/test", rl(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "allowed")
		})))

		server := &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadHeaderTimeout: 2 * time.Second,
		}

		log.Println("Middleware demo running at http://localhost:8080/test")
		log.Fatal(server.ListenAndServe())
}

