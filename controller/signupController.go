package controller

import (
	"Delfin/config"
	"Delfin/userModel"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Signup : regresa la vista con la forma
func Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dentro de Signup...")
	if userModel.AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", 303)
		return
	}
	config.TPL.ExecuteTemplate(w, "signup.html", nil)
}

// SignupProcess : post del signup
func SignupProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	newUser, err := userModel.CreateUser(r)
	if err != nil {
		//http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		http.Redirect(w, r, "/signup", 303)
		return
	}

	// crear sesion
	sID := uuid.Must(uuid.NewV4())
	cookie := &http.Cookie{
		Name:  "blablabla",
		Value: sID.String(),
		Path:  "/",
	}
	fmt.Println(newUser)

	// set cookie to writer
	http.SetCookie(w, cookie)

	// insertar valores a la tabla sessions
	userModel.InsertIntoSesion(cookie.Value, newUser.Email)

	// ejecutar template
	config.TPL.ExecuteTemplate(w, "home.html", newUser)
}
