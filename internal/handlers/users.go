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

func Users(res http.ResponseWriter, req *http.Request) {
	type newUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	reqUser := newUser{}
	err := decodeJSON(&reqUser, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSON(res, 500, s)
		return
	}

	hashPass, err := auth.HashPassword(reqUser.Password)
	if err != nil {
		s := fmt.Sprintf("`HashPassword()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, 500, s)
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
		returnJSON(res, 500, s)
		return
	}

	s := userJSON(&userRec)
	returnJSON(res, 201, s)
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	userRecs, err := cfg.DBQueries.GetUsers(req.Context())
	if err != nil {
		s := fmt.Sprintf("`GetUsers()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSON(res, 500, s)
	}

	stuff := []string{}
	for _, userRec := range userRecs {
		s := userJSON(&userRec)
		stuff = append(stuff, s)
	}
	s := fmt.Sprintf("[%s]", strings.Join(stuff, ","))
	returnJSON(res, 200, s)
}

// TODO: this is supposed to do marshalling stuff to map field names
// first pass just strings to make it work
func userJSON(user *database.User) string {
	// response chokes on white space, has to be ugly JSON

	id := fmt.Sprintf(`"%s": "%s"`, "id", user.ID)

	date := user.CreatedAt.Format(time.RFC3339)
	c := fmt.Sprintf(`"%s": "%s"`, "created_at", date)

	date = user.UpdatedAt.Format(time.RFC3339)
	u := fmt.Sprintf(`"%s": "%s"`, "updated_at", date)

	e := fmt.Sprintf(`"%s": "%s"`, "email", user.Email)

	s := fmt.Sprintf("{%s,%s,%s,%s}", id, c, u, e)

	// x := fmt.Sprintf("{\n\t%s,\n\t%s,\n\t%s,\n\t%s\n}\n", id, c, u, e)
	// fmt.Println(x)

	return s
}
