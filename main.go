package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/handlers"
)

func main() {
	fmt.Println("chirpy main()")
	cfg.DBinit()
	// cfg.TestDB()

	doStuff()
}

func doStuff() {
	godotenv.Load()

	const rootPath = "."
	const rootPrefix = "/app"
	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle(rootPrefix+"/",
		http.StripPrefix(rootPrefix,
			incHitCount(http.FileServer(http.Dir(rootPath))),
		),
	)

	mux.HandleFunc("GET /api/healthz", handlers.Healthz)
	mux.HandleFunc("POST /api/login", handlers.Login)

	mux.HandleFunc("POST /api/users", handlers.CreateUser)
	mux.HandleFunc("GET /api/users", handlers.GetUsers)

	mux.HandleFunc("POST /api/chirps", handlers.CreateChirp)
	mux.HandleFunc("GET /api/chirps", handlers.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.GetChirp)

	mux.HandleFunc("POST /api/refresh", handlers.Refresh)
	mux.HandleFunc("GET /api/refresh", handlers.GetRefresh)
	mux.HandleFunc("POST /api/revoke", handlers.Revoke)

	mux.HandleFunc("GET /admin/metrics", handlers.Metrics)
	mux.HandleFunc("POST /admin/reset", handlers.Reset)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func incHitCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
