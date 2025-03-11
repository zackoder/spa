package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user_id := models.GetUserId(cookie.Value)
	if user_id < 1 {
		fmt.Println("user not found")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	names := models.SelectUsers(user_id)

	json.NewEncoder(w).Encode(names)
}
