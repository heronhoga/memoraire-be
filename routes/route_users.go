package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
)

func UserRoutes(h *http.ServeMux) {
	h.HandleFunc("POST /register", handlers.Register)
}