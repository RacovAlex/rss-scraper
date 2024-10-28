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

// ApiConfig предоставляет конфигурацию для API-обработчиков,
// включая доступ к базе данных через объект Queries.
type ApiConfig struct {
	DB *database.Queries
}

// HandlerCreateUser создает нового пользователя, декодируя его имя из JSON-запроса.
// Регистрирует пользователя в базе данных и возвращает его данные.
func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := &parameters{}
	err := decoder.Decode(params)
	if err != nil {
		utils.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can't decode JSON: %v", err))
		return
	}

	// Создание нового пользователя в базе данных.
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Coudn't create user in: %v", err))
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.DatabaseUserToUser(user))
}

// HandlerGetUser возвращает информацию о текущем пользователе.
func (apiCfg *ApiConfig) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user))
}

// HandlerGetPostsForUsers возвращает список постов для текущего пользователя.
func (apiCfg *ApiConfig) HandlerGetPostsForUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUsers(r.Context(), database.GetPostsForUsersParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		utils.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get posts for user: %v", err))
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabasePostsToPosts(posts))
}
