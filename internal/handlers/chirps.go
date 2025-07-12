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

	isValid, userID := validateToken(res, req)
	if !isValid {
		return
	}

	c := chirpRequest{}
	err := decodeJSON(&c, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	if len(c.Body) > 140 {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Chirp is too long")
		returnJSONRes(res, http.StatusBadRequest, s)
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
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	s := dbChirpToJSON(&chirp)
	returnJSONRes(res, http.StatusCreated, s)
}

func GetChirps(res http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.DBQueries.GetChirps(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetChirps()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	stuff := []string{}
	for _, chirp := range chirps {
		s := dbChirpToJSON(&chirp)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSONRes(res, http.StatusOK, s)
}

func GetChirp(res http.ResponseWriter, req *http.Request) {

	isValid, id := validateRequestChirpID(res, req)
	if !isValid {
		return
	}

	isValid, chirp := getChirpRecord(res, req, id)
	if !isValid {
		return
	}

	s := dbChirpToJSON(&chirp)
	returnJSONRes(res, http.StatusOK, s)
}

func DeleteChirp(res http.ResponseWriter, req *http.Request) {
	isValid, id := validateRequestChirpID(res, req)
	if !isValid {
		return
	}

	isValid, userID := validateToken(res, req)
	if !isValid {
		return
	}

	isValid, chirp := getChirpRecord(res, req, id)
	if !isValid {
		return
	}

	if chirp.UserID != userID {
		returnJSONRes(res, http.StatusForbidden, "")
		return
	}

	err := cfg.DBQueries.DeleteChirp(req.Context(), id)
	if err != nil {
		s := fmt.Sprintf("`DeleteChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	returnJSONRes(res, http.StatusNoContent, "")
}

func getChirpRecord(res http.ResponseWriter, req *http.Request, id uuid.UUID) (bool, database.Chirp) {

	chirp, err := cfg.DBQueries.GetChirp(req.Context(), id)
	if err != nil {
		// TODO: is there a way to do this that does not suck
		if err.Error() == "sql: no rows in result set" {
			s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Couldn't get chirp")
			returnJSONRes(res, http.StatusNotFound, s)
			return false, chirp
		}

		s := fmt.Sprintf("`GetChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return false, chirp
	}

	return true, chirp
}

func validateRequestChirpID(res http.ResponseWriter, req *http.Request) (bool, uuid.UUID) {
	chirpID := req.PathValue("chirpID")
	// fmt.Printf("\nchirpID: %v\n", chirpID)

	id, err := uuid.Parse(chirpID)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "bad chirp ID")
		returnJSONRes(res, http.StatusBadRequest, s)
		return false, id
	}

	return true, id
}

func validateToken(res http.ResponseWriter, req *http.Request) (bool, uuid.UUID) {
	code := http.StatusUnauthorized

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("`GetBearerToken()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", http.StatusText(code))
		returnJSONRes(res, code, s)
		return false, uuid.Nil
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.JWTsecret)
	if err != nil {
		fmt.Println(err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", http.StatusText(code))
		returnJSONRes(res, code, s)
		return false, uuid.Nil
	}

	return true, userID
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
