package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Login(res http.ResponseWriter, req *http.Request) {
	type loginUser struct {
		Password      string `json:"password"`
		Email         string `json:"email"`
		ExpireSeconds int    `json:"expires_in_seconds"`
	}

	reqUser := loginUser{}
	err := decodeJSON(&reqUser, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	userRec, err := cfg.DBQueries.GetUserByEmail(req.Context(), reqUser.Email)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	err = auth.CheckPasswordHash(reqUser.Password, userRec.HashedPassword)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	expirationTime := time.Hour
	if reqUser.ExpireSeconds > 0 && reqUser.ExpireSeconds < cfg.TokenLimit {
		expirationTime = time.Duration(reqUser.ExpireSeconds) * time.Second
	}
	// fmt.Printf("\nexpirationTime: %v\n", expirationTime)

	tokenString, err := auth.MakeJWT(userRec.ID, cfg.Secret, expirationTime)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "create token failed")
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	s := userJSON(&userRec, tokenString)
	returnJSON(res, http.StatusOK, s)
}
