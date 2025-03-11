package controllers

import (
	"encoding/json"
	"net/http"

	"reat-time-forum/models"
)

func Servercategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories := models.SelectCategories(w)
	json.NewEncoder(w).Encode(categories)
}
