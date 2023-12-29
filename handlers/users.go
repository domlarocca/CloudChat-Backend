package handlers

import (
	"encoding/json"
	"net/http"
	"CloudChat/database"
	"log"
)


type RegistrationResponse struct {
	Message string `json:"message"`
}

type RegistrationRequest struct {
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var request RegistrationRequest
	log.Println(request)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Birthday == "" || request.Username == "" || request.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err = database.RegisterUser(request.Email, request.Birthday, request.Username, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := RegistrationResponse{Message: "User registered successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}