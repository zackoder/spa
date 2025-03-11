package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	postId, _ := strconv.Atoi(r.URL.Query().Get("id"))

	exits := models.CheckPost(postId)
	if !exits {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	dataComments := models.QueryComments(postId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataComments)
}
