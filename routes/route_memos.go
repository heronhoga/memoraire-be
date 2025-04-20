package routes

import (
	"net/http"

	"github.com/heronhoga/memoraire-be/handlers"
	"github.com/heronhoga/memoraire-be/utils"
)

func MemoRoutes(h *http.ServeMux) {
	//protected routes
	h.Handle("POST /memo/create", utils.CheckToken((http.HandlerFunc(handlers.CreateMemo))))
}