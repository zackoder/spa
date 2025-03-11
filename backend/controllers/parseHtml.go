package controllers

import (
	"net/http"

	"reat-time-forum/utils"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	utils.ParseAndExecute(w)
}
