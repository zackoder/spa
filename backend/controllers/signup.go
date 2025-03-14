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
		fmt.Println("error", err)
		utils.CreateJson(w, "there is an error try another time", http.StatusInternalServerError)
		return
	}
	user_id := models.Insert(w, userData.Nickname, userData.Firstname, userData.Lastname, userData.Gender, userData.Age, userData.Email, hashPassword)
	if user_id == -1 {
		return
	}
	
	w.WriteHeader(http.StatusCreated)

}


