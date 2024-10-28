package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"rss-scraper/internal/database"
	"rss-scraper/internal/models"
	"rss-scraper/pkg/utils"
)

// HandlerCreateFeedFollow создает подписку на канал (feed follow)
// для текущего пользователя, декодируя feed_id из JSON-запроса.
func (apiCfg *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can't decode JSON: %v", err))
		return
	}

	// Создает новую запись подписки на канал в базе данных.
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Can't create feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

// HandlerGetFeedFollows возвращает список всех подписок на каналы
// для текущего пользователя.
func (apiCfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Coudn't get feed follows: %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

// HandlerDeleteFeedFollow удаляет подписку на канал (feed follow)
// для текущего пользователя по заданному feedFollowID.
func (apiCfg *ApiConfig) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't parse feed follow ID: %v", err))
		return
	}

	// Удаляет подписку на канал в базе данных.
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete feed follow: %v", err))
	}

	utils.RespondWithJSON(w, http.StatusOK, struct{}{})
}
