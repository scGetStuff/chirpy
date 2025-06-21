package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("chirpy main()")
	doStuff()
}

func doStuff() {
	const rootPath = "."
	const rootPrefix = "/app"
	const port = "8080"
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// L5
	// mux.Handle("/", http.FileServer(http.Dir(".")))

	// L11
	mux.Handle(rootPrefix+"/",
		http.StripPrefix(rootPrefix, http.FileServer(http.Dir(rootPath))))
	mux.HandleFunc("/healthz", health)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func health(res http.ResponseWriter, req *http.Request) {
	// TODO: 503 Service Unavailable

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}
