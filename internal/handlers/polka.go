package handlers

import (
	"fmt"
	"net/http"

	"github.com/scGetStuff/chirpy/internal/auth"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

// TODO: lesson says request will never end untill 2xx code
// why did it have us return a 404 & 401?
// should 500 be changed to StatusNoContent?
func Polka(res http.ResponseWriter, req *http.Request) {
	type polkaData struct {
		UserID string `json:"user_id"`
	}
	type polkaRequest struct {
		Event string    `json:"event"`
		Data  polkaData `json:"data"`
	}

	reqStuff := polkaRequest{}
	err := decodeJSON(res, req, &reqStuff)
	if err != nil {
		return
	}

	reqKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		s := fmt.Sprintf("`GetAPIKey()` failed\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusUnauthorized, s)
		return
	}

	if reqKey != cfg.PolkaKey {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "key does not match")
		returnJSONResponse(res, http.StatusUnauthorized, s)
		return
	}

	if reqStuff.Event != "user.upgraded" {
		returnJSONResponse(res, http.StatusNoContent, "")
		return
	}

	userID, err := parseID(res, reqStuff.Data.UserID)
	if err != nil {
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
			returnJSONResponse(res, http.StatusNotFound, "")
			return
		}

		s := fmt.Sprintf("`UpdateUserRed()` failed:\n%v", err)
		s = fmt.Sprintf(`{"%s": "%s"}`, "error", s)
		returnJSONResponse(res, http.StatusInternalServerError, s)
		return
	}

	returnJSONResponse(res, http.StatusNoContent, "")
}
