package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func Validate_chirp(res http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	c := chirp{}
	err := decodeJSON(&c, req)
	if err != nil {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Something went wrong")
		returnJSON(res, 500, s)
		return
	}

	if len(c.Body) > 140 {
		s := fmt.Sprintf(`{"%s": "%s"}`, "error", "Chirp is too long")
		returnJSON(res, 400, s)
		return
	}

	c.Body = censor(c.Body)
	s := fmt.Sprintf(`{"%s": "%s"}`, "cleaned_body", c.Body)
	returnJSON(res, 200, s)
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
