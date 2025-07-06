package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Refresh(res http.ResponseWriter, req *http.Request) {

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("`GetBearerToken()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusBadRequest, s)
		return
	}

	userID, err := cfg.DBQueries.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		// TODO: is there a way to do this that does not suck
		if err.Error() == "sql: no rows in result set" {
			s := fmt.Sprintf(`{"%s": "%s"}`, "error", "refresh token either does not exist or is expired")
			returnJSON(res, http.StatusUnauthorized, s)
			return
		}

		s := fmt.Sprintf("`GetUserFromRefreshToken()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	expirationTime := time.Hour
	tokenAccess, err := auth.MakeJWT(userID, cfg.Secret, expirationTime)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "'MakeJWT()' failed")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	s := fmt.Sprintf(`{"%s": "%s"}`, "token", tokenAccess)
	returnJSON(res, http.StatusOK, s)
}

func GetRefresh(res http.ResponseWriter, req *http.Request) {

	tokenRecs, err := cfg.DBQueries.GetRefreshTokens(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetRefreshTokens()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	stuff := []string{}
	for _, tokenRec := range tokenRecs {
		data, err := json.Marshal(tokenRec)
		if err != nil {
			log.Fatal(err)
		}
		s := string(data)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))

	returnJSON(res, http.StatusOK, s)
}

func Revoke(res http.ResponseWriter, req *http.Request) {

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("`GetBearerToken()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusBadRequest, s)
		return
	}

	err = cfg.DBQueries.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		s := fmt.Sprintf("`GetRefreshTokens()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	res.WriteHeader(http.StatusNoContent)
}
