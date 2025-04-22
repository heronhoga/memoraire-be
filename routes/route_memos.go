package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func MemoRoutes(h *http.ServeMux) {
	//protected routes
	h.Handle("POST /memo/create", utils.CheckToken((http.HandlerFunc(handlers.CreateMemo))))
	h.Handle("GET /memo", utils.CheckToken((http.HandlerFunc(handlers.ReadMemo))))
	h.Handle("PUT /memo", utils.CheckToken((http.HandlerFunc(handlers.UpdateMemo))))
	h.Handle("DELETE /memo", utils.CheckToken((http.HandlerFunc(handlers.DeleteMemo))))
}