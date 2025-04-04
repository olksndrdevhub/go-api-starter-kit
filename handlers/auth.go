package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/oleksandrdevhub/go-api-starter-kit/db"
	"github.com/oleksandrdevhub/go-api-starter-kit/utils"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	UserId  int64  `json:"user_id,omitempty"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		http.Error(w, "First name and last name are required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// check if user already exists
	exists, err := db.CheckUserExistsByEmail(req.Email)
	if err != nil {
		http.Error(w, "Error checking user existence", http.StatusConflict)
		return
	}
	if exists {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// hash password
	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Error processing registration", http.StatusInternalServerError)
		return
	}

	// create user
	userID, err := db.CreateUser(req.Email, hashedPass, req.FirstName, req.LastName)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Error creationg user", http.StatusInternalServerError)
		return
	}

	// generate JWT token
	token, err := utils.GenerateJWTToken(userID, req.Email)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	// return success resp with token
	response := AuthResponse{
		Token:   token,
		Message: "Registration successful",
		UserId:  userID,
	}

	utils.WriteJson(w, http.StatusCreated, response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Invalid credentials: user not found", http.StatusUnauthorized)
		return
	}

	match, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil || !match {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Invalid credentials: password mismatch", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWTToken(user.ID, user.Email)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:   token,
		Message: "Login successful",
		UserId:  user.ID,
	}

	utils.WriteJson(w, http.StatusOK, response)
}
