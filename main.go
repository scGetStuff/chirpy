package main

import (
	"fmt"
	"log"
	"net/http"

	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/handlers"
)

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
			incHitCount(http.FileServer(http.Dir(rootPath))),
		),
	)
	mux.HandleFunc("GET /api/healthz", handlers.Healthz)

	// CH3 L4
	// mux.HandleFunc("GET /api/metrics", handlers.Metrics)
	mux.HandleFunc("GET /admin/metrics", handlers.Metrics)
	// mux.HandleFunc("POST /api/reset", handlers.Reset)
	mux.HandleFunc("POST /admin/reset", handlers.Reset)

	// CH4 L2
	mux.HandleFunc("POST /api/validate_chirp", handlers.Validate_chirp)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func incHitCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
