package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"rss-scraper/internal/database"
	"rss-scraper/internal/handlers"
	"rss-scraper/internal/models"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Инициализация конфигурации из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получение порта из переменных окружения
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("$DB_URL must be set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	db := database.New(conn)

	// TODO: Location?
	apiCfg := handlers.ApiConfig{
		DB: db,
	}

	go models.StartScraping(db, 10, time.Minute)

	// Инициализация маршрутизатора
	router := chi.NewRouter()

	// Настройка CORS для маршрутизатора с заданными параметрами.
	// Обеспечивает управление кросс-доменными запросами и доступом к ресурсам API.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Создаём маршрутизатор для версии API v1.
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)
	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandlerGetFeeds)

	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandlerGetPostsForUsers))

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))

	// Монтируем v1Router с префиксом /v1 на основной маршрутизатор для изоляции маршрутов и обеспечения версионирования API.
	router.Mount("/v1", v1Router)

	// Настройка HTTP-сервера с маршрутизатором и адресом
	srv := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Запуск сервера
	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
