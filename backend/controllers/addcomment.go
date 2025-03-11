package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Addcomment(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var comment utils.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		fmt.Println(err)
	}
	user_id := models.GetUserId(cookie.Value)
	if user_id == 0 {
		w.WriteHeader(401)
		return
	}
	models.InsertNewComment(user_id, comment)
	// db.Exec("INSERT INTO comments (user_id, post_id, comment, date) VALUES (?,?,?,strftime('%s', 'now'))", user_id, commet.PostId, commet.Comment)
}
