package routes

import (
	"log/slog"
	"net/http"

	"github.com/ibrahimijai/go-http-server/controllers"
	"github.com/ibrahimijai/go-http-server/middleware"
	"github.com/ibrahimijai/go-http-server/models"
)

// Router ties routing, middleware, and controllers together and is itself an
// http.Handler.
type Router struct {
	handler http.Handler
	logger  *slog.Logger
}

var _ http.Handler = (*Router)(nil)

func New(logger *slog.Logger, st models.URLStore) *Router {
	c := controllers.New(st, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", c.Health)
	mux.HandleFunc("POST /api/shorten", c.Shorten)
	mux.HandleFunc("GET /{code}", c.Redirect)

	rl := middleware.NewRateLimiter(5, 10) // 5 req/s steady, bursts of 10

	// Order reads outermost → innermost: RequestID first so everything
	// downstream has the ID, then Logging, then Recovery, then the rate limiter.
	chained := middleware.Chain(mux,
		middleware.RequestID,
		middleware.Logging(logger),
		middleware.Recovery(logger),
		rl.Middleware,
	)

	return &Router{handler: chained, logger: logger}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.handler.ServeHTTP(w, r)
}
