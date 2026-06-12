package controllers

import (
	"authentication/config"
	"authentication/helpers"
	db "authentication/internal/db"
	"authentication/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

func Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var req models.SignupRequest

		//Get user input
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		//Validate user input
		if validationErr := validate.Struct(req); validationErr != nil {
			helpers.WriteError(w, http.StatusBadRequest, validationErr.Error())
			return
		}

		count, err := config.Queries.CountByEmailOrPhone(ctx, db.CountByEmailOrPhoneParams{
			Email: req.Email,
			Phone: req.Phone,
		})
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if count > 0 {
			helpers.WriteError(w, http.StatusBadRequest, "Email or phone already exists")
			return
		}

		//Generate the rest of the user data
		userID := uuid.New()
		hashedPassword := helpers.HashPassword(req.Password)
		accessToken, refreshToken := helpers.GenerateTokens(req.Email, userID.String(), req.Role)

		_, err = config.Queries.CreateUser(ctx, db.CreateUserParams{
			UserID:       userID,
			FirstName:    req.First_name,
			LastName:     req.Last_name,
			Password:     hashedPassword,
			Email:        req.Email,
			Phone:        req.Phone,
			Role:         req.Role,
			Token:        &accessToken,
			RefreshToken: &refreshToken,
		})
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "User created successfully"})
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var req models.LoginRequest

		// Get user input
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		foundUser, err := config.Queries.GetUserByEmail(ctx, req.Email)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid email or password")
			return
		}

		//verify password
		passwordIsValid, _ := helpers.VerifyPassword(foundUser.Password, req.Password)
		if !passwordIsValid {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		token, refreshToken := helpers.GenerateTokens(foundUser.Email, foundUser.UserID.String(), foundUser.Role)
		if err := helpers.UpdateAllTokens(token, refreshToken, foundUser.UserID); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		foundUser.Password = "" // never expose the password hash

		helpers.WriteJSON(w, http.StatusOK, map[string]any{
			"user":          foundUser,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestedUserId := r.PathValue("id")

		// Get claims from the request context
		tokenClaims, ok := helpers.ClaimsFromContext(r.Context())
		if !ok {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if tokenClaims.Role != "ADMIN" && tokenClaims.UserID != requestedUserId {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		userID, err := uuid.Parse(requestedUserId)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid user id")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		user, err := config.Queries.GetUserByID(ctx, userID)
		if err != nil {
			helpers.WriteError(w, http.StatusNotFound, "User not found")
			return
		}

		user.Password = ""
		helpers.WriteJSON(w, http.StatusOK, user)
	}
}

func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve claims from the request context
		tokenClaims, ok := helpers.ClaimsFromContext(r.Context())
		if !ok {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Check if the user is an ADMIN
		if tokenClaims.Role != "ADMIN" {
			helpers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		users, err := config.Queries.ListUsers(ctx)
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		for i := range users {
			users[i].Password = "" // never expose password hashes
		}

		helpers.WriteJSON(w, http.StatusOK, users)
	}
}
