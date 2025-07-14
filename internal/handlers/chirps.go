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

	userID, err := validateToken(res, req)
	if err != nil {
		return
	}

	c := chirpRequest{}
	err = decodeJSON(res, req, &c)
	if err != nil {
		return
	}

	if len(c.Body) > 140 {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Chirp is too long")
		returnJSONResponse(res, http.StatusBadRequest, s)
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
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	s := dbChirpToJSON(&chirp)
	returnJSONResponse(res, http.StatusCreated, s)
}

func GetChirps(res http.ResponseWriter, req *http.Request) {
	var chirps []database.Chirp
	var dbErr error

	owner := req.URL.Query().Get("author_id")
	if owner != "" {
		userID, err := parseID(res, owner)
		if err != nil {
			return
		}
		chirps, dbErr = cfg.DBQueries.GetUserChirps(req.Context(), userID)
	} else {
		chirps, dbErr = cfg.DBQueries.GetChirps(req.Context())
	}
	if dbErr != nil {
		s := fmt.Sprintf("`GetChirps()` failed:\n%v", dbErr)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	stuff := []string{}
	for _, chirp := range chirps {
		s := dbChirpToJSON(&chirp)
		stuff = append(stuff, s)
	}

	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSONResponse(res, http.StatusOK, s)
}

func GetChirp(res http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	id, err := parseID(res, chirpID)
	if err != nil {
		return
	}

	chirp, err := getChirpRecord(res, req, id)
	if err != nil {
		return
	}

	s := dbChirpToJSON(&chirp)
	returnJSONResponse(res, http.StatusOK, s)
}

func DeleteChirp(res http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	id, err := parseID(res, chirpID)
	if err != nil {
		return
	}

	userID, err := validateToken(res, req)
	if err != nil {
		return
	}

	chirp, err := getChirpRecord(res, req, id)
	if err != nil {
		return
	}

	if chirp.UserID != userID {
		returnJSONResponse(res, http.StatusForbidden, "")
		return
	}

	err = cfg.DBQueries.DeleteChirp(req.Context(), id)
	if err != nil {
		s := fmt.Sprintf("`DeleteChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	returnJSONResponse(res, http.StatusNoContent, "")
}

func getChirpRecord(res http.ResponseWriter, req *http.Request, id uuid.UUID) (database.Chirp, error) {
	chirp, err := cfg.DBQueries.GetChirp(req.Context(), id)
	if err != nil {
		// TODO: is there a way to do this that does not suck
		if err.Error() == "sql: no rows in result set" {
			s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Couldn't get chirp")
			returnJSONResponse(res, http.StatusNotFound, s)
			return chirp, err
		}

		s := fmt.Sprintf("`GetChirp()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return chirp, err
	}

	return chirp, nil
}

func validateToken(res http.ResponseWriter, req *http.Request) (uuid.UUID, error) {
	code := http.StatusUnauthorized

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", err)
		returnJSONResponse(res, code, s)
		return uuid.Nil, err
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.JWTsecret)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", err)
		returnJSONResponse(res, code, s)
		return uuid.Nil, err
	}

	return userID, nil
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
