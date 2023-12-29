package handlers

import (
	"fmt"
	"log"
	"net/http"
	"CloudChat/database"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := database.DB.Ping()
	if err != nil {
		log.Println("Database Connection Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Database Connection Error 500")
		return
	}
	log.Println("Database Connection Good")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Database Connection OK")
}