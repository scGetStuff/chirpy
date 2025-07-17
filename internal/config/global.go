package config

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/scGetStuff/chirpy/internal/database"
)

// TODO: I did not like the Struct Method thing the lesson did
// the struct method being a handler function didn't feel right
// like it was mixing behavior that should be seperated

// unless a need arises for encapsulation, I'm just doing globals, because thats what they are

var FileServerHits = atomic.Int32{}
var DBQueries *database.Queries
var IsDev = false
var JWTsecret = ""
var TokenLimit = 60 * 60 // 1 hour in seconds
var PolkaKey = ""

func DBinit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error in `godotenv.Load()`:\n%v", err)
	}

	s := strings.ToLower(os.Getenv("PLATFORM"))
	// fmt.Printf("ENV: PLATFORM: %s\n", s)
	IsDev = (s == "dev")
	// fmt.Printf("IsDev: %v\n", IsDev)

	JWTsecret = os.Getenv("JWT_SECRET")
	// fmt.Printf("ENV: JWT_SECRET: %s\n", JWTsecret)

	PolkaKey = os.Getenv("POLKA_KEY")
	// fmt.Printf("ENV: POLKA_KEY: %s\n", PolkaKey)

	dbURL := os.Getenv("DB_URL")
	// fmt.Printf("ENV: DB_URL: %s\n", dbURL)
	if dbURL == "" {
		log.Fatal("error getting DB_URL from enviornment\n")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening DB:\n%v", err)
	}

	DBQueries = database.New(db)
	if DBQueries == nil {
		log.Fatalf("`DBQueries` bad stuff happened\n")
	}
}
