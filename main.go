package main

import (
	"Delfin/config"
	"Delfin/controller"
	"fmt"
	"log"
	"net/http"

	//libreria para mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cleanAll()
	http.HandleFunc("/", index)
	http.HandleFunc("/home", controller.Home)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/login/process", controller.LoginProcess)
	http.HandleFunc("/signup", controller.Signup)
	http.HandleFunc("/signup/process", controller.SignupProcess)
	http.HandleFunc("/logout", controller.Logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)

}

func index(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/home", http.StatusSeeOther)
	config.TPL.ExecuteTemplate(w, "home.html", nil)
}

func cleanAll() {
	fmt.Println("Deleting sessions...")
	_, err := config.DB.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		log.Fatalln(err)
	}
}
