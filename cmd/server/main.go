package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"rss-scraper/internal/database"
	"rss-scraper/internal/handlers"
	"rss-scraper/internal/models"

	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию из файла .env для настройки окружения
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем порт из переменных окружения, обязательно для запуска сервера
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	// Получаем URL базы данных из переменных окружения, обязательно для подключения к Postgres
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("$DB_URL must be set")
	}

	// Устанавливаем подключение к базе данных PostgreSQL
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	// Инициализируем новый объект Queries для взаимодействия с БД
	db := database.New(conn)

	// Создаем экземпляр ApiConfig для передачи зависимости базы данных в обработчики
	// TODO: Location?
	apiCfg := handlers.ApiConfig{
		DB: db,
	}

	// Запускаем процесс для регулярного парсинга RSS-лент с заданной частотой и числом горутин
	go models.StartScraping(db, 10, time.Minute)

	// Инициализация маршрутизатора
	router := chi.NewRouter()

	// Устанавливаем настройки CORS для маршрутизатора, разрешая доступ к API из указанных источников
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Создаем маршрутизатор для версии API v1, обеспечивая изоляцию маршрутов для удобного версионирования
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

	// Монтируем v1Router на основной маршрутизатор с префиксом /v1 для создания пространств имен и версионирования API
	router.Mount("/v1", v1Router)

	// Создаем и настраиваем HTTP-сервер, указывая маршрутизатор и порт для запуска
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
