package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("token")
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func CreateJson(w http.ResponseWriter, strErr string, codErr int) {
	err := ErrorResponse{Err: strErr, Code: codErr}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codErr)
	json.NewEncoder(w).Encode(err)
}

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}
