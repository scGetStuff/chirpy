package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func decodeJSON[T any](out *T, req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(out)

	return err
}

func returnJSON(res http.ResponseWriter, code int, json string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(json))
	if err != nil {
		log.Fatalf("returnJSON() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}

func returnTXT(res http.ResponseWriter, code int, msg string) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
	_, err := res.Write([]byte(msg))
	if err != nil {
		log.Fatalf("returnTXT() failed, that is not supposed to happen, I'm going to crash now\n%v", err)
	}
}
