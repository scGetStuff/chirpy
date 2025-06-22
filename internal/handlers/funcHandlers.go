package handlers

import (
	"fmt"
	"net/http"

	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Healthz(res http.ResponseWriter, req *http.Request) {
	// TODO: 503 Service Unavailable

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func Metrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	s := fmt.Sprintf("Hits: %d", cfg.FileServerHits.Load())
	res.Write([]byte(s))
}

func Reset(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	cfg.FileServerHits.Store(0)
	s := fmt.Sprintf("Reset: %d", cfg.FileServerHits.Load())
	res.Write([]byte(s))
}
