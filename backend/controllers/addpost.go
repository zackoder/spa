package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Addpost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": http.StatusText(http.StatusMethodNotAllowed)})
	}

	cookie := utils.CheckCookie(r)
	if cookie == nil {
		http.Redirect(w, r, "/signin", http.StatusUnauthorized)
		return
	}

	NewPost, user_id := models.SelectUserNicknameAndId(cookie.Value)

	json.NewDecoder(r.Body).Decode(&NewPost)
	message, postId := models.InsertPost(NewPost, user_id)
	if message != "" {
		fmt.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	NewPost.Id = postId
	json.NewEncoder(w).Encode(NewPost)
}
