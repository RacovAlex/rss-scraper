package main

import (
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"os"
	"rss-scraper/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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
