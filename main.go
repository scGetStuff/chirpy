package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
	"github.com/scGetStuff/chirpy/internal/handlers"
)

func main() {
	fmt.Println("chirpy main()")

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Printf("ENV: %s\n\n", dbURL)

	// doStuff()
	testDB(dbURL)
}

func testDB(dbURL string) {

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("open DB error")
	}

	dbQueries := database.New(db)
	_, err = dbQueries.CreateUser(context.Background(), "test@test.com")
	if err != nil {
		fmt.Printf("users failed: %w", err)
		log.Fatal("insert error")

	}

	users, err := dbQueries.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("users failed: %w", err)
		log.Fatal("select error")
	}

	fmt.Print(users)
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
