package handlers

import (
	"CloudChat/database"
	"encoding/json"
	"log"
	"net/http"
)

type LoginResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	log.Println(request)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	if request.Username == "" || request.Password == "" {
		http.Error(w, "All fields are required. Something was left blank!", http.StatusBadRequest)
		return
	}

	err = database.LoginUser(request.Username, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := LoginResponse{Message: "User logged in successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
