package controllers

import (
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Signout(w http.ResponseWriter, r *http.Request) {
	cookie := utils.CheckCookie(r)

	if cookie == nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	models.DeletingToken(cookie)

	http.SetCookie(w, &http.Cookie{
		Value:  "",
		Name:   "forum_token",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
