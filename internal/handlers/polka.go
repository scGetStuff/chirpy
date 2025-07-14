package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

// TODO: lesson saya request will never end untill 2xx code
// why did it have us return a 404?
// should 500 be changed to StatusNoContent?
func Polka(res http.ResponseWriter, req *http.Request) {
	type polkaData struct {
		UserID string `json:"user_id"`
	}
	type polkaRequest struct {
		Event string    `json:"event"`
		Data  polkaData `json:"data"`
	}

	reqKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		fmt.Printf("`GetAPIKey()` failed\n%v", err)
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", err.Error())
		returnJSONRes(res, http.StatusUnauthorized, s)
		return
	}
	if reqKey != cfg.PolkaKey {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "key does not match")
		returnJSONRes(res, http.StatusUnauthorized, s)
	}

	reqStuff := polkaRequest{}
	err = decodeJSON(&reqStuff, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	if reqStuff.Event != "user.upgraded" {
		returnJSONRes(res, http.StatusNoContent, "")
		return
	}

	userID, err := uuid.Parse(reqStuff.Data.UserID)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "bad user_id")
		returnJSONRes(res, http.StatusBadRequest, s)
		return
	}

	_, err = cfg.DBQueries.UpdateUserRed(req.Context(),
		database.UpdateUserRedParams{
			ID:          userID,
			IsChirpyRed: true,
		},
	)
	if err != nil {
		// TODO: is there a way to do this that does not suck
		if err.Error() == "sql: no rows in result set" {
			returnJSONRes(res, http.StatusNotFound, "")
			return
		}

		s := fmt.Sprintf("`UpdateUserRed()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONRes(res, http.StatusInternalServerError, s)
		return
	}

	returnJSONRes(res, http.StatusNoContent, "")
}
