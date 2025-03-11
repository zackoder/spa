package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if utils.CheckCookie(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		var req utils.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println(err)
		}
		if err := models.InsertUser(req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"message": "Somthing went wrong"})
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
	}
}
