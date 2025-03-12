package models

import (
	"fmt"
	"net/http"

	"reat-time-forum/utils"

	"golang.org/x/crypto/bcrypt"
)

func QueryComments(postId int) []utils.Comment {
	query := `SELECT c.comment, u.nickname, c.date
				FROM comments AS c 
				INNER JOIN users AS u ON u.id = c.user_id 
				WHERE c.post_id = ?;`
	rows, err := db.Query(query, postId)
	if err != nil {
		fmt.Println("quering comments", err)
		return nil
	}
	var dataComments []utils.Comment
	for rows.Next() {
		var comment utils.Comment
		comment.PostId = postId
		if err := rows.Scan(&comment.Comment, &comment.Username, &comment.CreationDate); err != nil {
			fmt.Println("error scan comments", err)
		}
		dataComments = append(dataComments, comment)
	}
	return dataComments
}

func SelectPosts(user_id int, offset string) []utils.Posts {
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
		return nil
	}
	var posts []utils.Posts

	for rows.Next() {
		var post utils.Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println(err)
			return nil
		}
		posts = append(posts, post)
	}

	defer rows.Close()

	GetPostsCategories(posts)
	return posts
}

func SelectUsers(user_id int) []utils.Name {
	var names []utils.Name

	query := `
	SELECT u.nickname
		FROM users u
		LEFT JOIN messages m
		    ON (u.id = m.sender_id OR u.id = m.reciever_id)
		    AND (m.sender_id = ? OR m.reciever_id = ?)
		WHERE u.id <> ?
		GROUP BY u.nickname
		ORDER BY 
    		MAX(m.creation_date) DESC,
    		u.nickname ASC;
	`
	rows, err := db.Query(query, user_id, user_id, user_id)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var name utils.Name
		if err := rows.Scan(&name.Name); err != nil {
			fmt.Println("error acrosed while scanning nicknmae :", err)
		}
		names = append(names, name)
	}

	defer rows.Close()
	return names
}

func SelectCategories(w http.ResponseWriter) []utils.Category {
	query := "SELECT name FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	var categories []utils.Category

	for rows.Next() {
		var category utils.Category
		if err := rows.Scan(&category.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		categories = append(categories, category)
	}
	defer rows.Close()
	return categories
}

func SelectUserNicknameAndId(token string) (utils.Posts, int) {
	var NewPost utils.Posts

	getUserId := `
   		SELECT u.nickname, s.user_id
    	FROM sessions s
    	JOIN users u ON u.id = s.user_id
    	WHERE s.token = ?;
		`

	var user_id int

	err := db.QueryRow(getUserId, token).Scan(&NewPost.Poster, &user_id)
	if err != nil {
		fmt.Println(err)
	}
	return NewPost, user_id
}

func Select(w http.ResponseWriter, username, password string) int {
	fmt.Println("username", username)
	// query := `SELECT id, password FROM users WHERE (nickname = ? AND password = ?) OR (email = ? AND password = ?)`
	query := `SELECT id, password FROM users WHERE nickname = ?  OR email = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println("err prepare", err)
		utils.CreateJson(w, "there is an error try another time", http.StatusInternalServerError)
		return -1
	}
	defer stmt.Close()
	var id int
	var hashPassword string

	err = stmt.QueryRow(username, username).Scan(&id, &hashPassword)
	// if err == sql.ErrNoRows{}
	if err != nil {
		fmt.Println("query row err", err)
		utils.CreateJson(w, "Username or Password not found", http.StatusMethodNotAllowed)
		return -1
	}
	if err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		fmt.Println("compare password err", err)
		utils.CreateJson(w, "your password incorrect", http.StatusMethodNotAllowed)
		return -1
	}
	fmt.Println("id", id)
	return id
}
