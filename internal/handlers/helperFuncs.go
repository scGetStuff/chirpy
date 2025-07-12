package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scGetStuff/chirpy/internal/database"
)

type userJSON struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func decodeJSON[T any](out *T, req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(out)

	return err
}

func returnJSONRes(res http.ResponseWriter, code int, json string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(json))
	if err != nil {
		log.Fatalf("returnJSON() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}

func returnTextRes(res http.ResponseWriter, code int, msg string) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(msg))
	if err != nil {
		log.Fatalf("returnTXT() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}

func dbUserToJSON(user *database.User) string {
	wraper := userJSON{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	return structToJSON(wraper)
}

func dbLoginToJSON(user *database.User, token string, refresh string) string {
	type stuff struct {
		userJSON
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	wraper := stuff{
		userJSON: userJSON{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: refresh,
	}

	return structToJSON(wraper)
}

func dbChirpToJSON(chirp *database.Chirp) string {
	type stuff struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		UserID    uuid.UUID `json:"user_id"`
		Body      string    `json:"body"`
	}

	wraper := stuff{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	}

	return structToJSON(wraper)
}

func structToJSON(stuff any) string {
	data, err := json.Marshal(stuff)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
		return ""
	}

	return string(data)
}
