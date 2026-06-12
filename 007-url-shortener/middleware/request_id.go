package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// ctxKey is an unexported context key type to avoid collisions with other
// packages.
type ctxKey int

const requestIDKey ctxKey = 0

// RequestID tags every request with an ID so logs across the stack correlate.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID") // honor upstream proxies
		if id == "" {
			id = newID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFrom retrieves the request ID from the context.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
