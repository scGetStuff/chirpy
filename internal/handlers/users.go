package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
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
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)
	returnJSON(res, http.StatusCreated, s)
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	userRecs, err := cfg.DBQueries.GetUsers(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetUsers()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
	}

	stuff := []string{}
	for _, userRec := range userRecs {
		s := dbUserToJSON(&userRec)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSON(res, http.StatusOK, s)
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
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("`GetBearerToken()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.Secret)
	if err != nil {
		fmt.Println(err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Unauthorized")
		returnJSON(res, http.StatusUnauthorized, s)
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
		returnJSON(res, http.StatusInternalServerError, s)
		return
	}

	s := dbUserToJSON(&userRec)

	returnJSON(res, http.StatusOK, s)
}

// TODO: still need to replace this
func userJSON(user *database.User, token string, refresh string) string {
	// response chokes on white space, has to be ugly JSON

	id := fmt.Sprintf(`"%s": "%s"`, "id", user.ID)

	date := user.CreatedAt.Format(time.RFC3339)
	cDate := fmt.Sprintf(`,"%s": "%s"`, "created_at", date)

	date = user.UpdatedAt.Format(time.RFC3339)
	uDate := fmt.Sprintf(`,"%s": "%s"`, "updated_at", date)

	email := fmt.Sprintf(`,"%s": "%s"`, "email", user.Email)

	t := ""
	if token != "" {
		t = fmt.Sprintf(`,"%s": "%s"`, "token", token)
	}

	r := ""
	if token != "" {
		r = fmt.Sprintf(`,"%s": "%s"`, "refresh_token", refresh)
	}

	s := fmt.Sprintf("{%s%s%s%s%s%s}", id, cDate, uDate, email, t, r)

	// fmt.Println(s)

	return s
}
