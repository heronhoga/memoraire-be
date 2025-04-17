package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func UserRoutes(h *http.ServeMux) {
	// Public routes
	h.Handle("POST /register", http.HandlerFunc(handlers.Register))
	h.Handle("POST /login", http.HandlerFunc(handlers.Login))

	// Protected route
	h.Handle("GET /logout", utils.CheckToken((http.HandlerFunc(handlers.Logout))))
}