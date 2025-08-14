package main

import (
	"log"
	"net/http"
)

func fileServer() http.Handler {
	return http.FileServer(http.Dir("./public"))
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /", fileServer())

	err := http.ListenAndServe(":1112", mux)
	if err != nil {
		log.Fatal(err)
	}
}
