package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func GetName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	// query := "SELECT nickname FROM users WHERE id = (SELECT user_id FROM sessions WHERE token = ?)"
	getId := models.GetUserId(cookie.Value)
	if getId < 1 {
		fmt.Println("user not found")
		w.WriteHeader(401)
		return
	}
	nickname := models.GetUserNickname(getId)

	// if err := db.QueryRow(query, cookie.Value).Scan(&nickname); err != nil {
	// 	fmt.Println(err)
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "unautorized"})
	// 	return
	// }

	json.NewEncoder(w).Encode(map[string]string{"nickname": nickname})
}
