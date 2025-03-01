package queries

import (
	"database/sql"
	"fmt"
	utils "reat-time-forum/structs"
)

func GetUserIdByNickname(db *sql.DB, nickname string) (int, error) {
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

func Getmessages(db *sql.DB, sender_id, reciever_id int, offset string) (error, []utils.Message) {

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

// func
