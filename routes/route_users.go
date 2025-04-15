package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
)

func UserRoutes(h *http.ServeMux) {
	h.HandleFunc("GET /register", handlers.Register)
}