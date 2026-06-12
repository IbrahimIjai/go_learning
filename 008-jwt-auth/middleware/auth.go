package middleware

import (
	"authentication/helpers"
	"log"
	"net/http"
	"strings"
)

// Authenticate is standard net/http middleware: it wraps a handler, validates
// the Bearer token, and stores the resulting claims in the request context.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		log.Printf("Raw Authorization Header: %s", authHeader) // Log the raw header

		if authHeader == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "authorization header required")
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token after trimming Bearer: %s", authHeader) // Log after trimming

		if authHeader == "" {
			helpers.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		claims, err := helpers.ValidateToken(authHeader)
		if err != nil {
			log.Printf("Token validation error: %v", err) // Log the validation error
			helpers.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Store claims in the request context for later use in controllers.
		ctx := helpers.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
