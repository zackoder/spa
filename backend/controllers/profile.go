package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	nickname := r.PathValue("nickname")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	user_id := models.GetUserId(cookie.Value)

	posts := models.FilterProfile(user_id, nickname, offset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
