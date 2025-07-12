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
	err := decodeJSON(&reqUser, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
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
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)
	returnJSONRes(res, http.StatusCreated, s)
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	userRecs, err := cfg.DBQueries.GetUsers(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetUsers()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	stuff := []string{}
	for _, userRec := range userRecs {
		s := dbUserToJSON(&userRec)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSONRes(res, http.StatusOK, s)
}

func PutUser(res http.ResponseWriter, req *http.Request) {
	type updateUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqUser := updateUser{}
	err := decodeJSON(&reqUser, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	isValid, userID := validateToken(res, req)
	if !isValid {
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
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)
	returnJSONRes(res, http.StatusOK, s)
}
