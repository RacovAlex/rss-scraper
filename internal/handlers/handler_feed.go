package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"rss-scraper/internal/database"
	"rss-scraper/internal/models"
	"rss-scraper/pkg/utils"
)

// HandlerCreateFeed создает новую запись канала (feed) в базе данных
// для текущего пользователя. Декодирует данные из JSON-запроса,
// создает feed и возвращает его в JSON-ответе.
func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can't decode JSON: %v", err))
		return
	}

	// Создание нового канала (feed) в базе данных.
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Coudn't create feed in: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedToFeed(feed))
}

// HandlerGetFeeds возвращает список всех доступных каналов (feeds).
func (apiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Coudn't get feeds: %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
}
