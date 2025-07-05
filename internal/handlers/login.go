package handlers

import (
	"fmt"
	"net/http"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
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
		returnJSON(res, 500, s)
		return
	}

	userRec, err := cfg.DBQueries.GetUserByEmail(req.Context(), reqUser.Email)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, 401, s)
		return
	}

	err = auth.CheckPasswordHash(reqUser.Password, userRec.HashedPassword)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, 401, s)
		return
	}

	s := userJSON(&userRec)
	returnJSON(res, 200, s)
}
