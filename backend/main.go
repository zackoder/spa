package main

import (
	"log"
	"net/http"

	"reat-time-forum/controllers"
	md "reat-time-forum/middleware"
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
	http.HandleFunc("/getNickName", md.Authorization(controllers.GetName))
	http.HandleFunc("/getusers", md.Authorization(controllers.GetUsers))
	http.HandleFunc("/sign-in", controllers.Signin)
	http.HandleFunc("/sign-up", controllers.Signup)
	http.HandleFunc("/signout", md.Authorization(controllers.Signout))
	http.HandleFunc("/addpost", md.Authorization(controllers.Addpost))
	http.HandleFunc("/posts", md.Authorization(controllers.GetPosts))
	http.HandleFunc("/comments", md.Authorization(controllers.GetComments))
	http.HandleFunc("/addcomment", md.Authorization(controllers.Addcomment))
	http.HandleFunc("/api/category/{categoryName}", md.Authorization(controllers.Handlecategories))
	http.HandleFunc("/api/{nickname}", md.Authorization(controllers.Profile))
	http.HandleFunc("/get_categories", md.Authorization(controllers.Servercategories))
	http.HandleFunc("/reactions", md.Authorization(controllers.HandleReactions))
	http.HandleFunc("/api/messages", md.Authorization(controllers.Fetchemessages))
	http.HandleFunc("/ws", md.Authorization(manager.ServeWebsocket))

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
