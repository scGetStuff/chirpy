package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	cfg "github.com/scGetStuff/chirpy/internal/config"
)

func Healthz(res http.ResponseWriter, req *http.Request) {
	// TODO: 503 Service Unavailable

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

// func Metrics(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
// 	res.WriteHeader(200)
// 	s := fmt.Sprintf("Hits: %d", cfg.FileServerHits.Load())
// 	res.Write([]byte(s))
// }

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

func Reset(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	cfg.FileServerHits.Store(0)
	s := fmt.Sprintf("Reset: %d", cfg.FileServerHits.Load())
	res.Write([]byte(s))
}

func Validate_chirp(res http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	decoder := json.NewDecoder(req.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		s := `{"error": "Something went wrong"}`
		res.WriteHeader(500)
		res.Write([]byte(s))
		return
	}

	if len(c.Body) > 140 {
		s := `{"error": "Chirp is too long"}`
		res.WriteHeader(400)
		res.Write([]byte(s))
		return
	}

	c.Body = censor(c.Body)

	res.WriteHeader(200)
	// s := `{"valid": true}`
	s := fmt.Sprintf(`{"cleaned_body": "%s"}`, c.Body)
	res.Write([]byte(s))
}

func censor(str string) string {
	badWords := []string{

		"kerfuffle",
		"sharbert",
		"fornax",
	}

	words := strings.Split(str, " ")

	for i, word := range words {
		word = strings.ToLower(word)
		for _, bad := range badWords {
			if word == bad {
				words[i] = "****"
			}
		}

	}

	return strings.Join(words, " ")
}
