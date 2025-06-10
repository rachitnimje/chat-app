package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware handles CORS for all requests
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment or use default
		allowedOrigins := []string{"http://localhost:3000"}
		if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
			allowedOrigins = strings.Split(envOrigins, ",")
		}

		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == strings.TrimSpace(allowedOrigin) {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartHTTPServer(router http.Handler, port string) error {
	// Wrap router with CORS middleware
	handler := CORSMiddleware(router)

	address := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(address, handler)
}
