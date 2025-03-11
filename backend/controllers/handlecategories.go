package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Handlecategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := utils.CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	user_id := models.GetUserId(cookie.Value)

	category, err := url.QueryUnescape(r.PathValue("categoryName"))
	if err != nil {
		fmt.Println(err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fmt.Println(err)
	}

	categoryId := models.GetCategory(category)
	if categoryId == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if categoryId == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	posts := models.FilterPostsByCategory(user_id, categoryId, offset)

	json.NewEncoder(w).Encode(posts)
}
