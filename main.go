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
	stuff := http.NewServeMux()
	server := http.Server{}
	server.Addr = ":8080"
	server.Handler = stuff

	stuff.Handle("/", http.FileServer(http.Dir("./pages")))

	server.ListenAndServe()
}
