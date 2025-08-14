package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
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

func BrotliMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "br") {
			w.Header().Set("Content-Encoding", "br")
			brWriter := brotli.NewWriter(w)
			defer brWriter.Close()
			bw := &brotliResponseWriter{ResponseWriter: w, Writer: brWriter}
			next.ServeHTTP(bw, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type brotliResponseWriter struct {
	http.ResponseWriter
	Writer *brotli.Writer
}

func (w *brotliResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
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
