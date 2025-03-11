package models

import (
	"database/sql"
	"fmt"
	"log"
)

func CheckPost(postId int) bool {
	checkpost := `
		SELECT EXISTS(
			SELECT 1 FROM posts WHERE id = ?
		);
	`
	var exists bool
	db.QueryRow(checkpost, postId).Scan(&exists)
	return exists
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
