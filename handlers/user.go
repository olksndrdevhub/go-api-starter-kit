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

func Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		log.Printf("ERORR: %v", errors.New("user id or email not found"))
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
