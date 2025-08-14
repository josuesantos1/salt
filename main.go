package main

import (
	"log"
	"net/http"
	"strings"
	"os"
	"github.com/andybalholm/brotli"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Address string `toml:"address"`
	Root    string `toml:"root"`
}

func loadConfig(path string) (*Config, error) {
	cfg := &Config{
		Address: "1112",
		Root:    "public",
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func fileServer(root string) http.Handler {
	return http.FileServer(http.Dir(root))
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
	cfg, err := loadConfig("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file not found, using defaults: %v", err)
		} else {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	log.Printf("Salt HTTP server starting on %s...", cfg.Address)

	mux := http.NewServeMux()
	mux.Handle("GET /", fileServer(cfg.Root))

	loggedMux := LogMiddleware(mux)

	err = http.ListenAndServe(":1112", loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}
