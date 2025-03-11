package models

import (
	"database/sql"
	"fmt"
	"strings"

	utils "reat-time-forum/utils"
)

var db *sql.DB

func GetUserIdByNickname(nickname string) (int, error) {
	getId := "SELECT id FROM users WHERE nickname = ?"

	stmt, err := db.Prepare(getId)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var user_id int
	if err := stmt.QueryRow(nickname).Scan(&user_id); err != nil {
		return 0, err
	}
	return user_id, nil
}

func Getmessages(sender_id, reciever_id int, offset string) (error, []utils.Message) {
	var msgs []utils.Message
	querytenmsg := `
		SELECT u.nickname, m.content, m.creation_date
		FROM messages m
		JOIN users u 
		ON u.id = reciever_id
		WHERE (sender_id = ? AND reciever_id = ?) 
		OR (sender_id = ? AND reciever_id = ?)
		ORDER BY m.id DESC 
		LIMIT 10 OFFSET ?;
	`
	rows, err := db.Query(querytenmsg, sender_id, reciever_id, reciever_id, sender_id, offset)
	if err != nil {
		fmt.Println("quering err:", err)
		return err, nil
	}

	defer rows.Close()

	for rows.Next() {
		var msg utils.Message
		err := rows.Scan(&msg.To, &msg.Content, &msg.CreatedAt)
		if err != nil {
			fmt.Println("scanning err:", err)
			return err, nil
		}
		msgs = append(msgs, msg)
	}
	return nil, msgs
}

func GetUserId(token string) int {
	var user_id int
	getUserId := "SELECT user_id FROM sessions WHERE token = ?"
	if err := db.QueryRow(getUserId, token).Scan(&user_id); err != nil {
		fmt.Println("getting user err:", err)
		return user_id
	}
	return user_id
}

func GetUserNickname(id int) string {
	var usernickname string
	getUserId := "SELECT nickname FROM users WHERE id = ?"
	if err := db.QueryRow(getUserId, id).Scan(&usernickname); err != nil {
		fmt.Println("getting user err:", err)
		return usernickname
	}
	return usernickname
}

func GetPostsCategories(posts []utils.Posts) {
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
