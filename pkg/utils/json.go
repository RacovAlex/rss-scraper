package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseWithError формирует JSON-ответ с сообщением об ошибке и заданным HTTP-кодом.
// Если код ошибки 5XX, сообщение логируется как критическое для дальнейшего анализа.
func ResponseWithError(w http.ResponseWriter, code int, message string) {
	// Логирует сообщение об ошибке для 5XX ошибок сервера.
	if code > 499 {
		log.Println("Responding with 5XX error: ", message)
	}

	// Определяет формат ответа с полем Error для единообразия.
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	// Вызывает функцию для отправки JSON-ответа с кодом ошибки.
	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

// RespondWithJSON отправляет JSON-ответ с данными и указанным статус-кодом.
// Устанавливает заголовок Content-Type и код статуса ответа.
// Если сериализация JSON завершается с ошибкой, возвращает HTTP 500.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	const fn = "utils.respondWithJSON"

	// Преобразует данные в JSON-формат.
	data, err := json.Marshal(payload)
	if err != nil {
		// Логирует и возвращает 500 Internal Server Error при ошибке преобразования.
		log.Printf("Failed to marshal json response: %v, error: %v", payload, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливает заголовок Content-Type на application/json.
	w.Header().Set("Content-Type", "application/json")
	// Устанавливает статусный код ответа.
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("%v: %v", fn, err)
	}
}
