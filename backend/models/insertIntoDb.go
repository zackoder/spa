package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"reat-time-forum/utils"
)

func Insertmsg(sender_id, receiver_id int, content string) error {
	query := "INSERT INTO messages (sender_id, reciever_id, content, creation_date) VALUES (?,?,?, strftime('%s', 'now'))"
	_, err := db.Exec(query, sender_id, receiver_id, content)
	return err
}

func InsertNewComment(user_id int, comment utils.Comment) {
	db.Exec("INSERT INTO comments (user_id, post_id, comment, date) VALUES (?,?,?,strftime('%s', 'now'))", user_id, comment.PostId, comment.Comment)
}

func InsertOrUpdateReaction(id int, target string, userId int, action string) (utils.Reactinos, error) {
	var reaction utils.Reactinos

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

func InsertPost(post utils.Posts, user_id int) (string, int) {
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

func InsertUser(user utils.SignupRequest) error {
	fmt.Println(user)
	query := "INSERT INTO users (nickname, first_name, last_name, age, gender, email, password) VALUES (?,?,?,?,?,?,?)"
	_, err := db.Exec(query, user.NickName, user.FirstName, user.LastName, user.Age, user.Gender, user.Email, user.Password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func InsertNewUser(userId int, uuid string) error {
	query := "INSERT INTO sessions (user_id, token) VALUES (?, ?)"

	_, err := db.Exec(query, userId, uuid)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func CreateSession(w http.ResponseWriter, userid int, session string) {
	query := `INSERT INTO sessions(user_id, token) 
			  VALUES(?, ?)
			  ON CONFLICT (token) 
			  DO UPDATE SET
				token = EXCLUDED.token, 
				creation_date = CURRENT_TIMESTAMP;
			`

	stmt, err := db.Prepare(query)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no rows")
		}
		fmt.Println("exec", err)
		utils.CreateJson(w, "there are an error try another time", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userid, session)
	if err != nil {
		fmt.Println("exec", err)
		utils.CreateJson(w, "there are an error try another time", http.StatusInternalServerError)
		return
	}
}

func Insert(w http.ResponseWriter, nickname string, firstname string, lastname string, gender string, age string, email string, password string) int {
	query := `INSERT INTO users(nickname, first_name, last_name, gender, age, email, password) VALUES(?, ?, ?, ?, ?, ?, ?);`
	stmt, err := db.Prepare(query)
	if err != nil {
		utils.CreateJson(w, "there is an error try another time", http.StatusInternalServerError)
		return -1
	}
	defer stmt.Close()

	res, err := stmt.Exec(nickname, firstname, lastname, gender, age, email, password)
	if err != nil {
		utils.CreateJson(w, "there is an error try another time", http.StatusInternalServerError)
		return -1
	}

	if err != nil && strings.Contains(err.Error(), "email") {
		utils.CreateJson(w, "email already used try another email", http.StatusFound)
		return -1
	}
	lastId, _ := res.LastInsertId()
	return int(lastId)
}
