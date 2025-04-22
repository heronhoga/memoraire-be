package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func UserRoutes(h *http.ServeMux) {
	// Public routes
	h.Handle("POST /register", utils.WithMiddleware(handlers.Register, utils.CheckKey))
	h.Handle("POST /login", utils.WithMiddleware(handlers.Login, utils.CheckKey))

	// Protected route
	h.Handle("GET /logout", utils.WithMiddleware(handlers.Logout, utils.CheckKey, utils.CheckToken))
}