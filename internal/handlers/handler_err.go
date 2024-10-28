package handlers

import (
	"net/http"

	"rss-scraper/pkg/utils"
)

// HandlerErr отправляет клиенту сообщение об ошибке с HTTP-кодом 400.
// Используется для тестирования обработки ошибок в маршрутах.
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithError(w, 400, "Something went wrong")
}
