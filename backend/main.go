package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	port := ":8080"

	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	http.HandleFunc("/", loging)
	http.ListenAndServe(port, nil)
}

func loging(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmp.Execute(w, nil)
}
