package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("chirpy main()")
	doStuff()
}

func doStuff() {
	mux := http.NewServeMux()
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = mux

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server.ListenAndServe()
}
