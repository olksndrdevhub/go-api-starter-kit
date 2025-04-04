package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/olksndrdevhub/go-api-starter-kit/db"
	"github.com/olksndrdevhub/go-api-starter-kit/middleware"
	"github.com/olksndrdevhub/go-api-starter-kit/utils"
)

type ProfileResponse struct {
	db.User
}

type ProfileUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ChangePasswordRequest struct {
	Password        string `json:"password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		log.Printf("ERORR: %v", errors.New("user id not found"))
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("ERROR: %v", err)
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"message": "user not found, bad credentials"})
		return
	}

	if r.Method == http.MethodGet {
		response := ProfileResponse{
			User: *user,
		}
		utils.WriteJson(w, http.StatusOK, response)
	} else if r.Method == http.MethodPatch {
		var req ProfileUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.FirstName) == "" && strings.TrimSpace(req.LastName) == "" {
			utils.WriteJson(w, http.StatusBadRequest, map[string]string{"message": "at least first name or last name are required"})
			return
		}

		updatedFirstName := user.FirstName
		if req.FirstName != "" {
			updatedFirstName = req.FirstName
		}
		updatedLastName := user.LastName
		if req.LastName != "" {
			updatedLastName = req.LastName
		}

		err = db.UpdateUser(userID, updatedFirstName, updatedLastName)
		if err != nil {
			log.Printf("ERROR: %v", err)
			utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
			return
		}

		utils.WriteJson(w, http.StatusOK, map[string]string{"message": "user updated"})
	} else if r.Method == http.MethodDelete {
		// TODO: delete user
		utils.WriteJson(w, http.StatusNoContent, nil)
	} else {
		utils.WriteJson(w, http.StatusMethodNotAllowed, map[string]string{"message": "method not allowed"})
	}
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		log.Printf("ERORR: %v", errors.New("user id not found"))
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	if strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.NewPassword) == "" || strings.TrimSpace(req.ConfirmPassword) == "" {
		http.Error(w, "Current password, new password and confirm password are required", http.StatusBadRequest)
		return
	}

	//check if current password is correct
	user, err := db.GetUserByID(userID)
	if err != nil {
		log.Printf("ERROR: %v", err)
		utils.WriteJson(w, http.StatusUnauthorized, map[string]string{"message": "user not found, bad credentials"})
		return
	}

	match, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// check if new password and confirm password match
	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "New password and confirm password do not match", http.StatusBadRequest)
		return
	}

	// check if new password is valid
	err = utils.ValidatePassword(req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// hash new password
	hashedNewPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// set new password
	err = db.ChangeUserPassword(userID, hashedNewPassword)
	if err != nil {
		log.Printf("ERROR: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "password changed"})
}
