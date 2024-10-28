package handlers

import (
	"net/http"

	"rss-scraper/pkg/utils"
)

// HandlerReadiness отвечает с HTTP 200 OK, сигнализируя, что сервер
// готов к работе. Обычно используется для проверки работоспособности.
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
