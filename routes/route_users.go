package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func UserRoutes(h *http.ServeMux) {
	// Public routes
	h.Handle("POST /register", utils.CheckKey(http.HandlerFunc(handlers.Register)))
	h.Handle("POST /login", utils.CheckKey(http.HandlerFunc(handlers.Login)))

	// Protected route
	h.Handle("GET /logout", utils.CheckKey(utils.CheckToken((http.HandlerFunc(handlers.Logout)))))
}