package controller

import (
	"Delfin/config"
	"Delfin/userModel"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Login : regresa la vista con la forma
func Login(w http.ResponseWriter, r *http.Request) {
	if userModel.AlreadyLoggedIn(r) {
		fmt.Println("Already logged in")
		http.Redirect(w, r, "/home", 303)
		return
	}
	config.TPL.ExecuteTemplate(w, "login.html", nil)
}

// LoginProcess : post del login
func LoginProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// Login
	u, err := userModel.Login(r)
	if err != nil {
		http.Redirect(w, r, "/login", 303)
		return
	}

	// crear sesion
	sID := uuid.Must(uuid.NewV4())
	cookie := http.Cookie{
		Name:  "blablabla",
		Value: sID.String(),
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	// insertar valores a la tabla sessions
	userModel.InsertIntoSesion(cookie.Value, u.Email)

	// ejecutar template
	config.TPL.ExecuteTemplate(w, "home.html", u)
}

// Logout : regresa la vista con la forma
func Logout(w http.ResponseWriter, r *http.Request) {
	if !userModel.AlreadyLoggedIn(r) {
		fmt.Println("Not logged int")
		http.Redirect(w, r, "/home", 303)
		return
	}

	// get cookie
	cookie, _ := r.Cookie("blablabla")

	// delete the session
	userModel.DeleteSesion(cookie.Value)

	// Kill cooke
	cookie = &http.Cookie{
		Name:   "blablabla",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

/*
  expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
    http.SetCookie(w, &cookie)
*/
