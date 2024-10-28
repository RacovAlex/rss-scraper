package handlers

import (
	"fmt"
	"net/http"

	"rss-scraper/internal/auth"
	"rss-scraper/internal/database"
	"rss-scraper/pkg/utils"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

// MiddlewareAuth проверяет наличие и валидность API-ключа,
// аутентифицирует пользователя, и передает его данные в обработчик.
// При ошибке аутентификации возвращает код ошибки 401 или 404.
func (apiCfg *ApiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("Couldn't get API key: %v", err))
		}

		// Получение пользователя по API-ключу для дальнейших действий.
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
