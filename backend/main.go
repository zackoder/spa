package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type SignupRequest struct {
	NickName  string `json:"nickName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Resp struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header.Get("Connection"))
	fmt.Println(r.Header.Get("Upgrade"))
	fmt.Println(r.Header.Get("Sec-WebSocket-Key"))
	fmt.Println(r.Header.Get("Sec-WebSocket-Version"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Connection", r.Header.Get("Connection"))
		fmt.Println("Upgrade", r.Header.Get("Upgrade"))
		fmt.Println("Sec-WebSocket-Key", r.Header.Get("Sec-WebSocket-Key"))
		fmt.Println("Sec-WebSocket-Version", r.Header.Get("Sec-WebSocket-Version"))
		fmt.Println("hello", err)
		return
	}
	defer conn.Close()
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Recieved message:", string(data))
		if err := conn.WriteMessage(messageType, data); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "my_db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertdb(db)

	port := ":8088"
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/signin", signinPage)
	http.HandleFunc("/sign-in", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/ws", handleConnection)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Page Not Found"})
		return
	}
	ParseAndExecute(w)
}

func ParseAndExecute(w http.ResponseWriter) {
	tmp, err := template.ParseFiles("../index.html")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal Server Error"})
		return
	}
	tmp.Execute(w, nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if CheckCookie(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		ParseAndExecute(w)
	}
	if r.Method == http.MethodPost {
		var req SignupRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Println(err)
		}
		insertUser(req)
	}
}

func insertUser(user SignupRequest) {
	query := "INSERT INTO users (nickname, first_name, last_name, age, gender, Email, password) VALUES (?,?,?,?,?,?,?)"
	_, err := db.Exec(query, user.NickName, user.FirstName, user.LastName, user.Age, user.Gender, user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		os.Exit(10)
	}
}

func signinPage(w http.ResponseWriter, r *http.Request) {
	if CheckCookie(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		ParseAndExecute(w)
	}
}

func signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := r.FormValue("email")
	password := r.FormValue("password")
	resp := CheckCredentials(email, password)
	fmt.Println(resp)
	if resp.Code == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func CheckCredentials(email, password string) Resp {
	var hashedPassword string
	query := "SELECT password FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return Resp{Message: "Email not found", Code: 0}
		}
		log.Println("Database error:", err)
		return Resp{Message: "Database error", Code: 0}
	}

	if hashedPassword == "" {
		return Resp{Message: "No password found for user", Code: 0}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return Resp{Message: "Incorrect Password", Code: 0}
	}
	return Resp{Message: "Login successful", Code: 1}
}

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}
	return cookie
}

func insertdb(db *sql.DB) {
	query := `
	  PRAGMA foreign_keys = ON;
	  CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname VARCHAR(50),
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		age VARCHAR(50),
		gender VARCHAR(20),
		email VARCHAR(100) UNIQUE,
		password VARCHAR(100)
	  );
	  CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER NOT NULL,
		token VARCHAR(255) UNIQUE,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, token),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	  );
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Error creating tables:", err)
	}
}
