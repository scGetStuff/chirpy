package handlers

import (
	"net/http"
)

func Healthz(res http.ResponseWriter, req *http.Request) {
	// TODO: 503 Service Unavailable

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}
