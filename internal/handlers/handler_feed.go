package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"rss-scraper/internal/database"
	"rss-scraper/internal/models"
	"rss-scraper/pkg/utils"
	"time"
)

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

func (apiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Coudn't get feeds: %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseFeedsToFeeds(feeds))
}
