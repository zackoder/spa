package controllers

//
import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"reat-time-forum/models"
	"reat-time-forum/utils"

	"github.com/gofrs/uuid"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("method Err")
		// w.WriteHeader(401)
		utils.CreateJson(w, "Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userLogin := utils.LoginUsers{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		fmt.Println("decode json err", err)
		// w.WriteHeader(500)
		utils.CreateJson(w, "Internal Server Error tstErr", http.StatusInternalServerError)
		return
	}

	fmt.Println(userLogin)

	// if strings.Contains(userLogin.Username, "@") {
	// 	if !utils.IsValidEmail(w, userLogin.Username) {
	// 		fmt.Println("valid email")
	// 		utils.CreateJson(w, "invalid email", http.StatusNotFound)
	// 		return
	// 	}
	// } else {
	// 	if !utils.IsValidName(w, userLogin.Username) {
	// 		fmt.Println("is valid nickname")
	// 		utils.CreateJson(w, "invalid nickname", http.StatusNotFound)
	// 		return
	// 	}
	// }
	// if !utils.IsValidPassword(w, userLogin.Password) {
	// 	fmt.Println("is valid password")
	// 	utils.CreateJson(w, "invalid password", http.StatusNotFound)
	// 	return
	// }

	id := models.Select(w, userLogin.Username, userLogin.Password)
	if id == -1 {
		utils.CreateJson(w, "invalid credentiels", http.StatusNotFound)
		return
	}

	uid, err := uuid.NewV4()
	if err != nil {
		fmt.Println("err uuid", err)
		utils.CreateJson(w, "try again", http.StatusInternalServerError)
		return
	}

	models.CreateSession(w, id, uid.String())

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    uid.String(),
		MaxAge:   int(time.Hour) * 24,
		HttpOnly: true,
		Path:     "/",
	})

	// w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w)
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func Signin(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var siginData utils.SigninRequest
// 	json.NewDecoder(r.Body).Decode(&siginData)
// 	message, userId := models.CheckCredentials(siginData.Userinpt, siginData.Password)
// 	if message != "" {
// 		w.WriteHeader(http.StatusNotFound)
// 		json.NewEncoder(w).Encode(map[string]string{"message": message})
// 		return
// 	}

// 	uuid, err := uuid.NewV7()
// 	if err != nil {
// 		http.Error(w, "try to sign in another time", http.StatusInternalServerError)
// 		return
// 	}

// 	cookie := http.Cookie{
// 		Name:     "forum_token",
// 		Value:    uuid.String(),
// 		Expires:  time.Now().Add(24 * time.Hour),
// 		HttpOnly: true,
// 	}
// 	if err := models.InsertNewUser(userId, uuid.String()); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	http.SetCookie(w, &cookie)
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
