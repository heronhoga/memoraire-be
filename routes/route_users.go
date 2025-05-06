package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func UserRoutes(h *http.ServeMux) {
	// Public routes
	h.Handle("POST /api/register", utils.WithMiddleware(handlers.Register, utils.CheckKey))
	h.Handle("POST /api/login", utils.WithMiddleware(handlers.Login, utils.CheckKey))

	// Protected route
	h.Handle("GET /api/logout", utils.WithMiddleware(handlers.Logout, utils.CheckKey, utils.CheckToken))
}