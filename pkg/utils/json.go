package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", message)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	const fn = "utils.respondWithJSON"
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal json response: %v, error: %v", payload, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("%v: %v", fn, err)
	}
}
