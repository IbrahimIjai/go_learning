package routes

import (
	"authentication/controllers"
	"authentication/middleware"
	"net/http"
)

// SetupRoutes registers all routes on a standard library ServeMux and returns
// it as an http.Handler. Method-based patterns and the {id} wildcard require
// Go 1.22+.
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.Handle("POST /signup", controllers.Signup())
	mux.Handle("POST /login", controllers.Login())

	// Protected routes (wrapped with the Authenticate middleware)
	mux.Handle("GET /users", middleware.Authenticate(controllers.GetUsers()))
	mux.Handle("GET /user/{id}", middleware.Authenticate(controllers.GetUser()))

	return mux
}
