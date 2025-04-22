package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func MemoRoutes(h *http.ServeMux) {
	//protected routes
	h.Handle("POST /memo/create", utils.CheckKey(utils.CheckToken((http.HandlerFunc(handlers.CreateMemo)))))
	h.Handle("GET /memo", utils.CheckKey(utils.CheckToken((http.HandlerFunc(handlers.ReadMemo)))))
	h.Handle("PUT /memo", utils.CheckKey(utils.CheckToken((http.HandlerFunc(handlers.UpdateMemo)))))
	h.Handle("DELETE /memo", utils.CheckKey(utils.CheckToken((http.HandlerFunc(handlers.DeleteMemo)))))
}