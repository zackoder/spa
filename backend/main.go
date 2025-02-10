package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
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
type Post struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection", r.Header.Get("Connection"))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
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

	port := ":8080"
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/getNickName", getName)
	http.HandleFunc("/signin", signinPage)
	http.HandleFunc("/sign-in", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/addpost", addpost)
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/category/{categoryName}", handlecategories)
	http.HandleFunc("/{nickname}", profile)
	http.HandleFunc("/ws", handleConnection)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

type Posts struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Poster    string   `json:"poster"`
	CreatedAt int      `json:"createdAt"`
	Reactions Reactios `json:"reactions"`
}

type Reactios struct {
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Action   string `json:"action"`
}

func getName(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	cookie := CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	query := "SELECT nickname FROM users WHERE id = (SELECT user_id FROM sessions WHERE token = ?)"
	var nickname string
	if err := db.QueryRow(query, cookie.Value).Scan(&nickname); err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(map[string]string{"message": "unautorized"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"nickname": nickname})
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if CheckCookie(r) == nil {
		fmt.Println("cookie err")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	_ = offset
	if err != nil {
		fmt.Println("parssing offset err:", err)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden"})
		return
	}

	var posts []Posts
	
	query := `
			SELECT 
			p.id,
			p.title,
			p.content,
			u.nickname,
			p.createdAt
			FROM posts p
			JOIN users u ON u.id = p.user_id
			ORDER BY p.id DESC 
			LIMIT 20 OFFSET ?;
	`
	rows, err := db.Query(query, 0)

	if err != nil {
		fmt.Println("quering err:", err)
		return
	}

	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt); err != nil {
			fmt.Println(err)
			return
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func handlecategories(w http.ResponseWriter, r *http.Request) {

}

func profile(w http.ResponseWriter, r *http.Request) {
	
}

func addpost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": http.StatusText(http.StatusMethodNotAllowed)})
	}

	cookie := CheckCookie(r)
	if cookie == nil {
		http.Redirect(w, r, "/signin", http.StatusUnauthorized)
		return
	}

	getUserId := "SELECT user_id FROM sessions WHERE token = ?"

	var user_id int

	err := db.QueryRow(getUserId, cookie.Value).Scan(&user_id)
	if err != nil {
		fmt.Println(err)
	}

	var NewPost Post
	json.NewDecoder(r.Body).Decode(&NewPost)
	if message := insertPost(NewPost, user_id); message != "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
	}

	fmt.Println(NewPost)
}

func insertPost(post Post, user_id int) string {
	if post.Title == "" {
		fmt.Println("title")
		return "Titel can not be empty"
	}

	if post.Content == "" {
		fmt.Println("content")
		return "Content can not be empty"
	}

	if len(post.Categories) == 0 {
		fmt.Println("categories")
		return "You need to choose at least one category"
	}

	query := "INSERT INTO posts (title, content, user_id, createdAt) VALUES (?,?, ?, strftime('%s', 'now'))"

	_, err := db.Exec(query, post.Title, post.Content, user_id)
	if err != nil {
		fmt.Println("inserting err:", err)
		return "Try to post another time"
	}
	return ""
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)
	if cookie == nil {
		http.Redirect(w, r, "/signin", http.StatusUnauthorized)
		return
	}
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
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
		if err := insertUser(req); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"message": "Somthing went wrong"})
			return
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
	}
}

func insertUser(user SignupRequest) error {
	fmt.Println(user)
	query := "INSERT INTO users (nickname, first_name, last_name, age, gender, email, password) VALUES (?,?,?,?,?,?,?)"
	_, err := db.Exec(query, user.NickName, user.FirstName, user.LastName, user.Age, user.Gender, user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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

type signinRequest struct {
	Userinpt string `json:"email"`
	Password string `json:"password"`
}

func signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var siginData signinRequest
	json.NewDecoder(r.Body).Decode(&siginData)
	message := CheckCredentials(siginData.Userinpt, siginData.Password)
	if message != "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	cookie := http.Cookie{
		Name:     "forum_token",
		Value:    "test",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	query := "INSERT INTO sessions (user_id, token) VALUES (?, ?)"
	_, err := db.Exec(query, 1, "test")
	if err != nil {
		fmt.Println(err)
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func CheckCredentials(email, password string) string {
	fmt.Println("email:", email, "password:", password)
	var hashedPassword string
	query := "SELECT password FROM users WHERE email = ? OR nickname = ?"
	err := db.QueryRow(query, email, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("here", err)
			return "Email not found"
		}
		log.Println("Database error:", err)
		fmt.Println(err)
		return "Database error"
	}

	// if hashedPassword == "" {
	// 	return "No password found for user"
	// }

	// if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
	// 	return "Incorrect Password"
	// }
	return ""
}

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("forum_token")
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

	  CREATE TABLE IF NOT EXISTS posts (
	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		createdAt INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	  );
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Error creating tables:", err)
	}
}
