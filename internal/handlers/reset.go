package handlers

import (
	"fmt"
	"log"
	"net/http"

	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Reset(res http.ResponseWriter, req *http.Request) {
	if !cfg.IsDev {
		code := http.StatusForbidden
		returnTextResponse(res, code, http.StatusText(code))
		return
	}

	err := cfg.DBQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Fatalf("`DeleteUsers()` failed:\n%v", err)
		return
	}

	cfg.FileServerHits.Store(0)
	s := fmt.Sprintf("Reset: %d", cfg.FileServerHits.Load())
	returnTextResponse(res, http.StatusOK, s)
}
