package middleware

import "net/http"

// Middleware wraps an http.Handler and returns a new one.
type Middleware func(http.Handler) http.Handler

// Chain composes middlewares so the first one listed is the outermost (runs
// first on the way in, last on the way out):
//
//	Chain(h, RequestID, Logging, Recovery) → RequestID(Logging(Recovery(h)))
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
