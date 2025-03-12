package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reat-time-forum/models"
	"reat-time-forum/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(401)
		fmt.Println("error methodPost")
		utils.CreateJson(w, "Method not Allowed test err", http.StatusMethodNotAllowed)
		return
	}

	// Get Data Elements
	// nickname := r.FormValue("nickname")
	// firstname := r.FormValue("firstname")
	// lastname := r.FormValue("lastname")
	// gender := r.FormValue("gender")
	// birthdate := r.FormValue("birthdate")
	// email := r.FormValue("email")
	// password := r.FormValue("password")
	// confirmpassword := r.FormValue("confirmpassword")
	var userData utils.Users
	json.NewDecoder(r.Body).Decode(&userData)
	if !utils.IsValidName(w, userData.Nickname) || !utils.IsValidName(w, userData.Firstname) || !utils.IsValidName(w, userData.Lastname) {
		return
	}
	fmt.Println(userData)
	if !utils.IsValidGender(w, userData.Gender) || !utils.IsValidAge(w, userData.Age) ||
		!utils.IsValidEmail(w, userData.Email) || !utils.IsValidPassword(w, userData.Password) ||
		!utils.IsValidCheckPassword(w, userData.Password, userData.Confirm_Password) {
		return
	}

	hashPassword, err := utils.HashPassword(userData.Password)
	if err != nil {
		utils.CreateJson(w, "there is an error try another time", http.StatusInternalServerError)
		return
	}
	user_id := models.Insert(w, userData.Nickname, userData.Firstname, userData.Lastname, userData.Gender, userData.Age, userData.Email, hashPassword)
	if user_id == -1 {
		return
	}

	// uid, err := uuid.NewV4()
	// if err != nil {
	// 	utils.CreateJson(w, "there is an error try another time 2", http.StatusInternalServerError)
	// 	return
	// }
	// utils.CreateSession(w, user_id, uid.String())

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "token",
	// 	Value:    uid.String(),
	// 	MaxAge:   int(time.Hour) * 24,
	// 	HttpOnly: true,
	// 	Path:     "/",
	// })
	// http.Redirect(w, r, "/login", http.StatusSeeOther)
	w.WriteHeader(http.StatusCreated)

	// fmt.Println("data submit")
	// json.NewEncoder(w).Encode(map[string]string{"data": "success"})
}

// func Signup(w http.ResponseWriter, r *http.Request) {
// 	if utils.CheckCookie(r) != nil {
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		var req utils.SignupRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			fmt.Println(err)
// 		}
// 		if err := models.InsertUser(req); err != nil {
// 			json.NewEncoder(w).Encode(map[string]string{"message": "Somthing went wrong"})
// 			return
// 		}
// 		http.Redirect(w, r, "/signin", http.StatusFound)
// 	}
// }
