package main

import (
	"log"
	"net/http"

	"reat-time-forum/controllers"
	"reat-time-forum/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	models.Insertdb()

	port := ":8080"
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	manager := controllers.NewManager()

	http.HandleFunc("/", controllers.HomePage)
	http.HandleFunc("/getNickName", controllers.GetName)
	http.HandleFunc("/getusers", controllers.GetUsers)
	http.HandleFunc("/sign-in", controllers.Signin)
	http.HandleFunc("/sign-up", controllers.Signup)
	http.HandleFunc("/signout", controllers.Signout)
	http.HandleFunc("/addpost", controllers.Addpost)
	http.HandleFunc("/posts", controllers.GetPosts)
	http.HandleFunc("/comments", controllers.GetComments)
	http.HandleFunc("/addcomment", controllers.Addcomment)
	http.HandleFunc("/api/category/{categoryName}", controllers.Handlecategories)
	http.HandleFunc("/api/{nickname}", controllers.Profile)
	http.HandleFunc("/get_categories", controllers.Servercategories)
	http.HandleFunc("/reactions", controllers.HandleReactions)
	http.HandleFunc("/api/messages", controllers.Fetchemessages)
	http.HandleFunc("/ws", manager.ServeWebsocket)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
