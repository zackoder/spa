package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (m *Manager) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user_id := getUserId(cookie.Value)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := NewClient(conn, m, user_id)
	m.addClient(client)
	go client.readmessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client]; ok {
		client.Connection.Close()
		delete(m.clients, client)
	}
}

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

	manager := NewManager()

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/getNickName", getName)
	http.HandleFunc("/getusers", getUsers)
	http.HandleFunc("/signin", signinPage)
	http.HandleFunc("/sign-in", signin)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signout", signout)
	http.HandleFunc("/addpost", addpost)
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/api/category/{categoryName}", handlecategories)
	http.HandleFunc("/api/{nickname}", profile)
	http.HandleFunc("/get_categories", servercategories)
	http.HandleFunc("/reactions", handleReactions)
	http.HandleFunc("/ws", manager.serveWebsocket)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

type Name struct {
	Name string `json:"nickname"`
}

type ClientList map[*Client]bool

type Client struct {
	Connection *websocket.Conn
	manager    *Manager
	Client_id  int
}

type Message struct {
	To      string `json:"to"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (c *Client) readmessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		messagetype, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
			}
			break
		}
		var msg Message

		if err := json.Unmarshal(payload, &msg); err != nil {
			fmt.Println(err)
		}
		query := "SELECT id FROM users WHERE nickname = ?"

		var receiver_id int
		if err := db.QueryRow(query, msg.To).Scan(&receiver_id); err != nil {
			return
		}
		fmt.Println(receiver_id)
		fmt.Println("type", messagetype)
		fmt.Println(msg)
		fmt.Println(c.Client_id)
		if !isOnlien(c, receiver_id) {
			fmt.Println("receiver is offline", msg.To)
			c.Connection.WriteMessage(websocket.TextMessage, []byte(`{"status": "failed", "message": "User is offline"}`))
			continue
		}
	}
}

func isOnlien(c *Client, receiver_id int) bool {
	for client := range c.manager.clients {
		if client.Client_id == receiver_id {
			return true
		}
	}
	return false
}

func NewClient(conn *websocket.Conn, manager *Manager, user_id int) *Client {
	return &Client{
		Connection: conn,
		manager:    manager,
		Client_id:  user_id,
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cooki := CheckCookie(r)
	if cooki == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var names []Name

	user_id := getUserId(cooki.Value)
	query := "SELECT nickname FROM users WHERE id <> ?"
	rows, err := db.Query(query, user_id)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var name Name
		if err := rows.Scan(&name.Name); err != nil {
			fmt.Println("error acrosed while scanning nicknmae :", err)
		}
		names = append(names, name)
	}

	defer rows.Close()
	json.NewEncoder(w).Encode(names)
}

type Category struct {
	Name string `json:"name"`
}

func handleReactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie := CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	target := r.URL.Query().Get("target")
	action := r.URL.Query().Get("action")
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid post/comment ID", http.StatusBadRequest)
		return
	}

	if action != "like" && action != "dislike" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	reaction, err := insertOrUpdateReaction(id, target+"_id", cookie.Value, action)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println("DB error:", err)
		return
	}

	json.NewEncoder(w).Encode(reaction)
}

func insertOrUpdateReaction(id int, target, token, action string) (Reactinos, error) {
	var reaction Reactinos

	var userId int

	queryId := "SELECT user_id FROM sessions WHERE token = ?"
	if err := db.QueryRow(queryId, token).Scan(&userId); err != nil {
		return reaction, err
	}

	var existingReaction string
	checkQuery := "SELECT reaction_type FROM reactions WHERE user_id = ? AND " + target + " = ?"
	err := db.QueryRow(checkQuery, userId, id).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO reactions (user_id, "+target+", reaction_type) VALUES (?, ?, ?)", userId, id, action)
		if err != nil {
			fmt.Println("inserting err", err)
			return reaction, err
		}
	} else if err == nil {
		if existingReaction == action {
			_, err = db.Exec("DELETE FROM reactions WHERE user_id = ? AND "+target+" = ?", userId, id)
			if err != nil {
				fmt.Println("deleting err")
			}
		} else {
			_, err = db.Exec("UPDATE reactions SET reaction_type = ? WHERE user_id = ? AND "+target+" = ?", action, userId, id)
			if err != nil {
				fmt.Println("deleting err")
			}

		}
		if err != nil {
			return reaction, err
		}
	} else {
		return reaction, err
	}

	countQuery := `
		SELECT
			(SELECT COUNT(*) FROM reactions WHERE ` + target + ` = ? AND reaction_type = 'like') AS likes,
			(SELECT COUNT(*) FROM reactions WHERE ` + target + ` = ? AND reaction_type = 'dislike') AS dislikes
	`
	err = db.QueryRow(countQuery, id, id).Scan(&reaction.Likes, &reaction.Dislikes)
	if err != nil {
		fmt.Println("selecting likes err", err)
		return reaction, err
	}
	reaction.Action = action
	return reaction, nil
}

func servercategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var categories []Category

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}
	defer rows.Close()
	json.NewEncoder(w).Encode(categories)
}

func signout(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)

	if cookie == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		// http.Redirect(w, r, "/signin", http.StatusUnauthorized)
		return
	}

	if _, err := db.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value); err != nil {
		fmt.Println(err)
	}

	http.SetCookie(w, &http.Cookie{
		Value:  "",
		Name:   "forum_token",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/signin", http.StatusOK)
}

type Posts struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Poster     string    `json:"poster"`
	CreatedAt  int       `json:"createdAt"`
	Categories []string  `json:"categories"`
	Reactions  Reactinos `json:"reactions"`
}

type Reactinos struct {
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
	cookie := CheckCookie(r)
	if cookie == nil {
		fmt.Println("cookie err")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}

	offset := r.URL.Query().Get("offset")

	var posts []Posts
	user_id := getUserId(cookie.Value)

	query := `
			SELECT 
			p.id,
			p.title,
			p.content,
			u.nickname,
			p.createdAt,
			COALESCE(SUM(CASE WHEN r.reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN r.reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes,
			COALESCE((
				SELECT reaction_type
				FROM reactions
				WHERE user_id = ? AND post_id = p.id
				), '') AS user_reaction
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN reactions r ON r.post_id = p.id
			GROUP BY p.id
			ORDER BY p.id DESC
			LIMIT 20 OFFSET ?;
	`

	rows, err := db.Query(query, user_id, offset)

	if err != nil {
		fmt.Println("quering err:", err)
		return
	}

	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println(err)
			return
		}
		posts = append(posts, post)
	}

	defer rows.Close()

	getPostsCategories(posts)

	json.NewEncoder(w).Encode(posts)
}

func getUserId(token string) int {
	var user_id int
	getUserId := "SELECT user_id FROM sessions WHERE token = ?"
	if err := db.QueryRow(getUserId, token).Scan(&user_id); err != nil {
		fmt.Println("getting user err:", err)
		return user_id
	}
	return user_id
}

func getPostsCategories(posts []Posts) {
	if len(posts) > 0 {
		placeholders := strings.Repeat("?,", len(posts)-1) + "?"
		getCategoriesQuery := `
        SELECT pc.post_id, c.name
        FROM posts_categories pc
        JOIN categories c ON pc.category_id = c.id
        WHERE pc.post_id IN (` + placeholders + `)
    `

		args := make([]interface{}, len(posts))
		for i, post := range posts {
			args[i] = post.Id
		}

		categoryRows, err := db.Query(getCategoriesQuery, args...)
		if err != nil {
			fmt.Println("Error fetching categories:", err)
			return
		}
		defer categoryRows.Close()

		categoryMap := make(map[int][]string)
		for categoryRows.Next() {
			var postID int
			var categoryName string
			if err := categoryRows.Scan(&postID, &categoryName); err != nil {
				fmt.Println("Error scanning categories:", err)
				return
			}
			categoryMap[postID] = append(categoryMap[postID], categoryName)
		}

		for i := range posts {
			if categories, exists := categoryMap[posts[i].Id]; exists {
				posts[i].Categories = categories
			}
		}
	}
}

func handlecategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := CheckCookie(r)
	if cookie == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	var user_id int
	getUserId := "SELECT user_id FROM sessions WHERE token = ?;"

	if err := db.QueryRow(getUserId, cookie.Value).Scan(&user_id); err != nil {
		fmt.Println("getting user id inside categories function", err)
		return
	}
	var categoryId int
	category, err := url.QueryUnescape(r.PathValue("categoryName"))

	if err != nil {
		fmt.Println(err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fmt.Println(err)
	}

	getCategoryIdQuery := "SELECT id FROM categories WHERE name = ?"
	if err := db.QueryRow(getCategoryIdQuery, category).Scan(&categoryId); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("there is no category named " + category)
			http.Error(w, "there is no category named "+category, http.StatusNotFound)
			return
		} else {
			fmt.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	getPostsQuery := `
		SELECT 
			p.id, 
			p.title, 
			p.content,
			u.nickname,
			p.createdAt,
			COALESCE(COUNT(CASE WHEN r.reaction_type = 'like' THEN 1 END), 0) AS likes,
			COALESCE(COUNT(CASE WHEN r.reaction_type = 'dislike' THEN 1 END), 0) AS dislikes,
			COALESCE((
				SELECT reaction_type
				FROM reactions
				WHERE user_id = ? AND post_id = p.id
			), '') AS user_reaction
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN posts_categories pc ON p.id = pc.post_id
		LEFT JOIN reactions r ON r.post_id = p.id
		WHERE pc.category_id = ?
		GROUP BY p.id, p.title, p.content, u.nickname, p.createdAt
		ORDER BY p.id DESC
		LIMIT 20 OFFSET ?;
	`

	rows, err := db.Query(getPostsQuery, user_id, categoryId, offset)
	if err != nil {
		fmt.Println("erro while getting posts:", err)
	}

	defer rows.Close()

	var posts []Posts

	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println("Row scan error:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	defer rows.Close()

	getPostsCategories(posts)

	json.NewEncoder(w).Encode(posts)
}

func profile(w http.ResponseWriter, r *http.Request) {
	cookie := CheckCookie(r)
	if cookie == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	nickname := r.PathValue("nickname")
	offset := r.URL.Query().Get("offset")
	user_id := getUserId(cookie.Value)
	fmt.Println("nickname:", nickname)
	fmt.Println("offset:", offset)
	fmt.Println("user id:", user_id)
	query := `
			SELECT
			    p.id,
			    p.title,
			    p.content,
			    u.nickname,
			    p.createdAt,
			    COALESCE(SUM(CASE WHEN r.reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
			    COALESCE(SUM(CASE WHEN r.reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes,
			    COALESCE((
			        SELECT reaction_type
			        FROM reactions
			        WHERE user_id = ? AND post_id = p.id
			    ), '') AS user_reaction
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN reactions r ON r.post_id = p.id
			WHERE u.nickname = ?
			GROUP BY p.id
			ORDER BY p.id DESC
			LIMIT 20 OFFSET ?;
	`
	rows, err := db.Query(query, user_id, nickname, offset)
	if err != nil {
		fmt.Printf("error while geting user %s profile\n", nickname)
	}

	var posts []Posts

	for rows.Next() {
		var post Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println(err)
			return
		}
		posts = append(posts, post)
	}

	getPostsCategories(posts)

	defer rows.Close()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
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
	var NewPost Posts

	getUserId := `
   		SELECT u.nickname, s.user_id
    	FROM sessions s
    	JOIN users u ON u.id = s.user_id
    	WHERE s.token = ?;
		`

	var user_id int

	err := db.QueryRow(getUserId, cookie.Value).Scan(&NewPost.Poster, &user_id)
	if err != nil {
		fmt.Println(err)
	}

	json.NewDecoder(r.Body).Decode(&NewPost)
	message, postId := insertPost(NewPost, user_id)
	if message != "" {
		fmt.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	NewPost.Id = postId
	json.NewEncoder(w).Encode(NewPost)
}

func insertPost(post Posts, user_id int) (string, int) {
	if post.Title == "" {
		fmt.Println("title")
		return "Titel can not be empty", 0
	}

	if post.Content == "" {
		fmt.Println("content")
		return "Content can not be empty", 0
	}

	if len(post.Categories) == 0 {
		fmt.Println("categories")
		return "You need to choose at least one category", 0
	}
	query := "INSERT INTO posts (title, content, user_id, createdAt) VALUES (?,?, ?, strftime('%s', 'now'))"

	res, err := db.Exec(query, post.Title, post.Content, user_id)

	if err != nil {
		fmt.Println("inserting err:", err)
		return "Try to post another time", 0
	}

	lastId, _ := res.LastInsertId()

	for _, category := range post.Categories {
		getcategoryId := "SELECT id FROM categories WHERE name = ?"
		var categoryId int
		if err := db.QueryRow(getcategoryId, category).Scan(&categoryId); err != nil {
			if err == sql.ErrNoRows {
				return "there is no catergory named " + category, 0
			} else {
				fmt.Println("selecting and enserting categories", err.Error())
				return "Inertnal server err", 0
			}
		}

		insertPostCategory := "INSERT INTO posts_categories (post_id, category_id) VALUES (? , ?)"
		db.Exec(insertPostCategory, lastId, categoryId)
	}

	return "", int(lastId)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	// db.Exec("DELETE * FROM sessions")
	ParseAndExecute(w)
	// if r.URL.Path != "/" {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(map[string]string{"message": "Page Not Found"})
	// 	return
	// }
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
	ParseAndExecute(w)
}

type signinRequest struct {
	Userinpt string `json:"email"`
	Password string `json:"password"`
}

func signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var siginData signinRequest
	json.NewDecoder(r.Body).Decode(&siginData)
	fmt.Println(siginData)
	message, userId := CheckCredentials(siginData.Userinpt, siginData.Password)
	if message != "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "try to sign in another time", http.StatusInternalServerError)
		return
	}

	// fmt.Println(uuid)

	cookie := http.Cookie{
		Name:     "forum_token",
		Value:    uuid.String(),
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	query := "INSERT INTO sessions (user_id, token) VALUES (?, ?)"

	_, err = db.Exec(query, userId, uuid.String())
	if err != nil {
		fmt.Println(err)
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func CheckCredentials(email, password string) (string, int) {

	var hashedPassword string
	var userId int

	query := "SELECT id, password FROM users WHERE email = ? OR nickname = ?"
	err := db.QueryRow(query, email, email).Scan(&userId, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("here", err)
			return "Email not found", 0
		}
		log.Println("Database error:", err)
		return "Database error", 0
	}

	return "", userId
}

func CheckCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("forum_token")
	if err != nil {
		fmt.Println(err)
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
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	  );

	  CREATE TABLE IF NOT EXISTS categories (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name varchar(255) NOT NULL UNIQUE
	  );

	  CREATE TABLE IF NOT EXISTS comments (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		comment TEXT NOT NULL,
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	  );

	  CREATE TABLE IF NOT EXISTS reactions (
   		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		user_id INTEGER NOT NULL,
   		post_id INTEGER,
   		comment_id INTEGER,
    	reaction_type TEXT NOT NULL,
    	date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    	FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    	FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    	CHECK (
        	(post_id IS NOT NULL AND comment_id IS NULL) OR
        	(comment_id IS NOT NULL AND post_id IS NULL)
    	)
	  );

	  CREATE TABLE IF NOT EXISTS posts_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id,category_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	  );
	  
	  CREATE TABLE IF NOT EXISTS messages (
	  	id INTEGER PRIMARY KEY AUTOINCREMENT,
	  	sender_id INTEGER NOT NULL,
		reciever_id INTEGER NOT NULL,
		content TEXT,
		is_read BOOLEAN DEFAULT FALSE,
		creation_date INTEGER,
		FOREIGN KEY (sender_id) REFERENCES users(id),
		FOREIGN KEY (reciever_id) REFERENCES users(id)
	  );
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Error creating tables:", err)
	}
	db.Exec(`INSERT INTO categories (name) VALUES
('Technology'),
('Science'),
('Health & Wellness'),
('Business & Finance'),
('Education'),
('Entertainment'),
('Sports'),
('Politics'),
('Travel'),
('Lifestyle'),
('Artificial Intelligence'),
('Cybersecurity'),
('Software Development'),
('Blockchain & Cryptocurrency'),
('Gadgets & Reviews'),
('Web Development'),
('Cloud Computing'),
('Gaming'),
('Nutrition'),
('Mental Health'),
('Exercise & Fitness'),
('Meditation & Mindfulness'),
('Medical News'),
('Movies & TV Shows'),
('Music'),
('Literature & Books'),
('Photography'),
('Fashion'),
('Food & Cooking'),
('History'),
('DIY & Crafts');
`)
}
