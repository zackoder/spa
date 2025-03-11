package models

import (
	"fmt"
	"net/http"
)

func DeletingToken(cookie *http.Cookie) {
	if _, err := db.Exec("DELETE FROM sessions WHERE token = ? AND EXISTS (SELECT 1 FROM sessions WHERE token = ?);", cookie.Value, cookie.Value); err != nil {
		fmt.Println(err)
	}
}
