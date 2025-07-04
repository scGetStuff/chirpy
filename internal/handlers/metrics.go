package handlers

import (
	"fmt"
	"net/http"

	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Metrics(res http.ResponseWriter, req *http.Request) {
	s := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	s = fmt.Sprintf(s, cfg.FileServerHits.Load())

	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte(s))
}
