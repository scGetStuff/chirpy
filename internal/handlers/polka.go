package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	cfg "github.com/scGetStuff/chirpy/internal/config"
	"github.com/scGetStuff/chirpy/internal/database"
)

func Polka(res http.ResponseWriter, req *http.Request) {
	type polkaData struct {
		UserID string `json:"user_id"`
	}
	type polkaRequest struct {
		Event string    `json:"event"`
		Data  polkaData `json:"data"`
	}

	reqStuff := polkaRequest{}
	err := decodeJSON(&reqStuff, req)
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
