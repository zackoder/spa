package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "my_db.db")
	if err != nil {
		fmt.Println(err)
	}

	insertdb(db)
	port := ":8080"
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signin", signin)
	http.HandleFunc("/signup", signup)

	if err := (http.ListenAndServe(port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(`{message: Page Not Found}`)
		return
	}

	ParsAndExcut(w)
}

func ParsAndExcut(w http.ResponseWriter) {
	tmp, err := template.ParseFiles("../index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)
	if cookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("singup")
		return
	}
	ParsAndExcut(w)
}

func signin(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)
	if cookie != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("singin")
		return
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println(email, password)
	}
	if r.Method == http.MethodGet {
		ParsAndExcut(w)
	}
}

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}
	return cookie
}

func insertdb(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname VARCHAR(50),
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		age VARCHAR(50),
		gender VARCHAR(20),
		Email VARCHAR(100),
		password VARCHAR(100)
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}
