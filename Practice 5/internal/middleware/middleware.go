package middleware

import (
	"log"
	"net/http"
	"time"
)

var ValidAPIKey = "mysecretkey"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("timestamp=%s method=%s endpoint=%s", start.Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" || apiKey != ValidAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "401 Unauthorized: missing or invalid X-API-KEY"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
