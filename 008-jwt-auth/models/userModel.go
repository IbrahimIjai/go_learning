package models

// SignupRequest is the JSON body accepted by POST /signup.
// Validation tags are processed by go-playground/validator.
type SignupRequest struct {
	First_name string `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  string `json:"last_name" validate:"required,min=2,max=100"`
	Password   string `json:"password" validate:"required,min=6"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	Role       string `json:"role" validate:"required,oneof=ADMIN USER"`
}

// LoginRequest is the JSON body accepted by POST /login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
