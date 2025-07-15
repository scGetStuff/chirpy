package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {
	type newUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqUser := newUser{}
	err := decodeJSON(res, req, &reqUser)
	if err != nil {
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	userRec, err := cfg.DBQueries.CreateUser(req.Context(),
		database.CreateUserParams{
			Email:          reqUser.Email,
			HashedPassword: hashPass,
		},
	)
	if err != nil {
		s := fmt.Sprintf("`CreateUser()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)
	returnJSONResponse(res, http.StatusCreated, s)
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	if !cfg.IsDev {
		code := http.StatusForbidden
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", http.StatusText(code))
		returnJSONResponse(res, code, s)
		return
	}

	userRecs, err := cfg.DBQueries.GetUsers(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetUsers()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	stuff := []string{}
	for _, userRec := range userRecs {
		s := dbUserToJSON(&userRec)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSONResponse(res, http.StatusOK, s)
}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	type updateUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqUser := updateUser{}
	err := decodeJSON(res, req, &reqUser)
	if err != nil {
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	userID, err := validateToken(res, req)
	if err != nil {
		return
	}

	userRec, err := cfg.DBQueries.UpdateUser(req.Context(),
		database.UpdateUserParams{
			ID:             userID,
			Email:          reqUser.Email,
			HashedPassword: hashPass,
		},
	)
	if err != nil {
		s := fmt.Sprintf("`UpdateUser()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)
	returnJSONResponse(res, http.StatusOK, s)
}
