package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("forum_token")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cookie
}

func ParseAndExecute(w http.ResponseWriter) {
	tmp, err := template.ParseFiles("../template/index.html")
	if err != nil {
		fmt.Println("error parse", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}
	tmp.Execute(w, nil)
}
