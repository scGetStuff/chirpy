package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
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

func DBinit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error in `godotenv.Load()`:\n%v", err)
	}

	dbURL := os.Getenv("DB_URL")
	fmt.Printf("ENV: %s\n\n", dbURL)
	if dbURL == "" {
		log.Fatalf("error getting DB_URL from enviornment\n")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening DB:\n%v", err)
	}

	DBQueries = database.New(db)
	if DBQueries == nil {
		log.Fatalf("this is not supposed to happen\n")
	}
}

func TestDB() {
	_, err := DBQueries.CreateUser(context.Background(), "test@test.com")
	if err != nil {
		fmt.Printf("`CreateUser()` failed: \n%v", err)
		log.Fatal("insert error")

	}

	users, err := DBQueries.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("`GetUsers()` failed: \n%v", err)
		log.Fatal("select error")
	}

	for _, user := range users {
		printUser(user)
	}
}

func printUser(user database.User) {
	fmt.Println()
	fmt.Printf("ID:    %v\n", user.ID)
	fmt.Printf("Email: %v\n", user.Email)
}
