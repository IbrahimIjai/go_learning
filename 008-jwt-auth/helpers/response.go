package helpers

import (
	"context"
	"encoding/json"
	"net/http"
)

// WriteJSON writes the given payload as a JSON response with the given status.
// It replaces gin's c.JSON.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError is a small convenience wrapper for error responses.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// contextKey is an unexported type to avoid collisions in request contexts.
type contextKey string

const claimsKey contextKey = "claims"

// ContextWithClaims stores the authenticated user's claims in the request context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves the claims previously stored by the auth middleware.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}
