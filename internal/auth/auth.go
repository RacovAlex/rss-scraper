package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey извлекает API Key из заголовка HTTP запроса.
// Пример:
// Authorization: ApiKey {insert apikey here}
func GetAPIKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("authorization header is missing")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("authorization header is invalid")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("first part of authorization header is invalid")
	}
	return vals[1], nil
}
