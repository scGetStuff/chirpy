package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

func Login(res http.ResponseWriter, req *http.Request) {
	type loginUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqUser := loginUser{}
	err := decodeJSON(&reqUser, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	userRec, err := cfg.DBQueries.GetUserByEmail(req.Context(), reqUser.Email)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSONRes(res, http.StatusUnauthorized, s)
		return
	}

	err = auth.CheckPasswordHash(reqUser.Password, userRec.HashedPassword)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSONRes(res, http.StatusUnauthorized, s)
		return
	}

	expirationTime := time.Hour
	tokenAccess, err := auth.MakeJWT(userRec.ID, cfg.Secret, expirationTime)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "'MakeJWT()' failed")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	tokenRefresh := auth.MakeRefreshToken()
	_, err = cfg.DBQueries.CreateRefreshToken(req.Context(),
		database.CreateRefreshTokenParams{
			Token:  tokenRefresh,
			UserID: userRec.ID,
		},
	)
	if err != nil {
		s := fmt.Sprintf("`CreateRefreshToken()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	s := dbLoginToJSON(&userRec, tokenAccess, tokenRefresh)
	returnJSONRes(res, http.StatusOK, s)
}
