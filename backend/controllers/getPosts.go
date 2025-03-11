package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		fmt.Println("cookie err")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	offset := r.URL.Query().Get("offset")

	user_id := models.GetUserId(cookie.Value)

	posts := models.SelectPosts(user_id, offset)

	json.NewEncoder(w).Encode(posts)
}
