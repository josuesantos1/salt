package main

import (
	"log"
	"net/http"
)

func fileServer() http.Handler {
	return http.FileServer(http.Dir("./public"))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(recorder, r)
		log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL.Path, recorder.status)
	})
}

func main() {

	log.Println("Salt HTTP server starting on :1112...")

	mux := http.NewServeMux()

	mux.Handle("GET /", fileServer())

	loggedMux := LogMiddleware(mux)

	err := http.ListenAndServe(":1112", loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}
