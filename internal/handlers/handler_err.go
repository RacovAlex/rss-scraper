package handlers

import (
	"net/http"
	"rss-scraper/pkg/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithError(w, 400, "Something went wrong")
}
