package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/scGetStuff/chirpy/internal/database"
)

type userStructForJSON struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsRed     bool      `json:"is_chirpy_red"`
}

func decodeJSON[T any](res http.ResponseWriter, req *http.Request, out *T) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(out)

	if err != nil {
		s := fmt.Sprintf("`decodeJSON()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
	}

	return err
}

func returnJSONResponse(res http.ResponseWriter, code int, json string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(json))
	if err != nil {
		log.Fatalf("returnJSONResponse() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}

func returnTextResponse(res http.ResponseWriter, code int, msg string) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(msg))
	if err != nil {
		log.Fatalf("returnTextResponse() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}

func dbUserWrapper(user *database.User) userStructForJSON {
	return userStructForJSON{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsRed:     user.IsChirpyRed,
	}
}

func dbUserToJSON(user *database.User) string {
	return structToJSON(dbUserWrapper(user))
}

func dbLoginToJSON(user *database.User, token string, refresh string) string {
	type stuff struct {
		userStructForJSON
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	wraper := stuff{
		userStructForJSON: dbUserWrapper(user),
		Token:             token,
		RefreshToken:      refresh,
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

func parseID(res http.ResponseWriter, sID string) (uuid.UUID, error) {
	uID, err := uuid.Parse(sID)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "bad UUID string")
		returnJSONResponse(res, http.StatusBadRequest, s)
		return uID, err
	}

	return uID, nil
}
