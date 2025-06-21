package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

var apiCfg = apiConfig{
	fileserverHits: atomic.Int32{},
}

func main() {
	fmt.Println("chirpy main()")
	doStuff()
}

func doStuff() {
	const rootPath = "."
	const rootPrefix = "/app"
	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// L5
	// mux.Handle("/", http.FileServer(http.Dir(".")))

	// L11
	mux.Handle(rootPrefix+"/",
		http.StripPrefix(rootPrefix,
			apiCfg.middlewareMetricsInc(
				http.FileServer(http.Dir(rootPath)))))
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/metrics", count)
	mux.HandleFunc("/reset", reset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func health(res http.ResponseWriter, req *http.Request) {
	// TODO: 503 Service Unavailable

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func count(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	s := fmt.Sprintf("Hits: %d", apiCfg.fileserverHits.Load())
	res.Write([]byte(s))
}

func reset(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	apiCfg.fileserverHits.Store(0)
	s := fmt.Sprintf("reset: %d", apiCfg.fileserverHits.Load())
	res.Write([]byte(s))
}
