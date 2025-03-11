package controllers

import (
	"encoding/json"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Fetchemessages(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)

	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sender_id := models.GetUserId(cookie.Value)
	offset := r.URL.Query().Get("offset")

	receiverNickname := r.URL.Query().Get("to")

	receiver_id, err := models.GetUserIdByNickname(receiverNickname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err, Messages := models.Getmessages(sender_id, receiver_id, offset)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(Messages)
}
