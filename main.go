package CloudChat_Backend

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the chat application!")
	})

	r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
		// Handle chat-related logic
	})

	// Add more routes as needed

	http.ListenAndServe(":8080", r)
}
