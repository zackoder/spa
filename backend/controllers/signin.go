package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"reat-time-forum/models"
	"reat-time-forum/utils"

	"github.com/gofrs/uuid"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var siginData utils.SigninRequest
	json.NewDecoder(r.Body).Decode(&siginData)
	message, userId := models.CheckCredentials(siginData.Userinpt, siginData.Password)
	if message != "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "try to sign in another time", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "forum_token",
		Value:    uuid.String(),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	if err := models.InsertNewUser(userId, uuid.String()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
