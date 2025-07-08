package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

func CreateChirp(res http.ResponseWriter, req *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("`GetBearerToken()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.Secret)
	if err != nil {
		fmt.Println(err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	c := chirpRequest{}
	err = decodeJSON(&c, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	if len(c.Body) > 140 {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Chirp is too long")
		returnJSON(res, http.StatusBadRequest, s)
		return
	}
	c.Body = censor(c.Body)

	chirp, err := cfg.DBQueries.CreateChirp(req.Context(),
		database.CreateChirpParams{
			Body:   c.Body,
			UserID: userID,
		},
	)
	if err != nil {
		s := fmt.Sprintf("`CreateChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	s := dbChirpToJSON(&chirp)
	returnJSON(res, http.StatusCreated, s)
}

func GetChirps(res http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.DBQueries.GetChirps(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetChirps()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	stuff := []string{}
	for _, chirp := range chirps {
		s := dbChirpToJSON(&chirp)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSON(res, http.StatusOK, s)
}

func GetChirp(res http.ResponseWriter, req *http.Request) {

	chirpID := req.PathValue("chirpID")
	// fmt.Printf("\nchirpID: %v\n", chirpID)

	id, err := uuid.Parse(chirpID)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "bad chirp ID")
		returnJSON(res, http.StatusBadRequest, s)
		return
	}

	chirp, err := cfg.DBQueries.GetChirp(req.Context(), id)
	if err != nil {
		// TODO: is there a way to do this that does not suck
		if err.Error() == "sql: no rows in result set" {
			s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Couldn't get chirp")
			returnJSON(res, http.StatusNotFound, s)
			return
		}

		s := fmt.Sprintf("`GetChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	s := dbChirpToJSON(&chirp)
	returnJSON(res, http.StatusOK, s)
}

func censor(str string) string {
	badWords := []string{

		"kerfuffle",
		"sharbert",
		"fornax",
	}

	words := strings.Split(str, " ")

	for i, word := range words {
		word = strings.ToLower(word)
		for _, bad := range badWords {
			if word == bad {
				words[i] = "****"
			}
		}

	}

	return strings.Join(words, " ")
}
