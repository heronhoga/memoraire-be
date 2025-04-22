package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func MemoRoutes(h *http.ServeMux) {
	//protected routes
	h.Handle("POST /memo/create", utils.WithMiddleware(handlers.CreateMemo, utils.CheckKey, utils.CheckToken))
	h.Handle("GET /memo", utils.WithMiddleware(handlers.ReadMemo, utils.CheckKey, utils.CheckToken))
	h.Handle("PUT /memo", utils.WithMiddleware(handlers.UpdateMemo, utils.CheckKey, utils.CheckToken))
	h.Handle("DELETE /memo", utils.WithMiddleware(handlers.DeleteMemo, utils.CheckKey, utils.CheckToken))
}