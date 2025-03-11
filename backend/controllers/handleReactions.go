package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func HandleReactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	userId := models.GetUserId(cookie.Value)

	target := r.URL.Query().Get("target")
	action := r.URL.Query().Get("action")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post/comment ID", http.StatusBadRequest)
		return
	}

	if action != "like" && action != "dislike" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	reaction, err := models.InsertOrUpdateReaction(id, target+"_id", userId, action)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("DB error:", err)
		return
	}

	json.NewEncoder(w).Encode(reaction)
}
