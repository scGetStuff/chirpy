package handlers

import (
	"fmt"
	"log"
	"net/http"

	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Reset(res http.ResponseWriter, req *http.Request) {

	if !cfg.IsDev {
		returnTXTRes(res, 403, "Forbidden")
		return
	}

	err := cfg.DBQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Fatalf("`DeleteUsers()` failed:\n%v", err)
		return
	}

	cfg.FileServerHits.Store(0)
	s := fmt.Sprintf("Reset: %d", cfg.FileServerHits.Load())
	returnTXTRes(res, 200, s)
}
