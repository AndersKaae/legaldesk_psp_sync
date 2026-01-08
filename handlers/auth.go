package handlers

import (
	"log"
	"net/http"
	"crypto/subtle" // For constant-time comparison
)

// RequireBasicAuth is a middleware that enforces HTTP Basic Authentication.
func RequireBasicAuth(next http.HandlerFunc, requiredUser, requiredPass string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(requiredUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(requiredPass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Unauthorized access attempt from %s", r.RemoteAddr)
			return
		}

		// If authentication is successful, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
