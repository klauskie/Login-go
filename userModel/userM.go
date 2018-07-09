package userModel

import (
	"Delfin/config"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Email     string
	Firstname string
	Lastname  string
	Password  string
}

// CreateUser :
func CreateUser(r *http.Request) (user, error) {

	newUser := user{}
	newUser.Email = r.FormValue("email")
	newUser.Firstname = r.FormValue("firstname")
	newUser.Lastname = r.FormValue("lastname")
	p := r.FormValue("password")
	p2 := r.FormValue("password2")

	// validar nulos
	if newUser.Email == "" || newUser.Firstname == "" || newUser.Lastname == "" || p == "" {
		fmt.Println("Error en validar nulos...")
		return newUser, errors.New("400. Bad request. All fields must be complete")
	}

	// validar si es unico
	if !findExisting(newUser.Email) {
		fmt.Println("Error en validar si es unico...")
		return newUser, errors.New("400. Bad request. Email allready taken")
	}

	// contraseñas coinciden
	if p != p2 {
		fmt.Println("Contraseñas no son iguales...")
		return newUser, errors.New("Internal server error. (Contraseñas son iguales)")
	}

	// encriptar contraseña
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		//http.Error(w, "Internal server error. (Pedo con encriptar)", http.StatusInternalServerError)
		fmt.Println("Error en encriptar la contraseña...")
		return newUser, errors.New("Internal server error. (Pedo con encriptar)")
	}

	// insertar user
	_, err = config.DB.Exec("INSERT INTO usuario(email, fname, lname, password) VALUES (?, ?, ?, ?)", newUser.Email, newUser.Firstname, newUser.Lastname, bs)
	if err != nil {
		fmt.Println("Error en insertar valores a la base de datos...")
		return newUser, errors.New("500. Internal Server Error. (Insertar Usuario)" + err.Error())
	}

	// Todo bien
	return newUser, nil
}

// Login : Se encarga de trabajar con el cliente y la base de datos
func Login(r *http.Request) (user, error) {

	un := r.FormValue("email")
	p := r.FormValue("password")

	u := user{}

	row := config.DB.QueryRow("SELECT * FROM usuario WHERE email = ?", un)
	var bs string
	err := row.Scan(&u.Email, &u.Firstname, &u.Lastname, &bs)
	if err != nil {
		fmt.Println("Error: email no registrado")
		return u, errors.New("Internal server error. (Pedo con validacion)")
	}
	err = bcrypt.CompareHashAndPassword([]byte(bs), []byte(p))
	if err != nil {
		fmt.Println("Error: email y/o contraseña no coinciden")
		return user{}, errors.New("Internal server error. (Pedo con validacion)")
	}
	// Validar usuario y contraseña con la Base de datos
	u.Password = bs
	return u, nil
}

// *----------------Funciones Auxiliares-------------*

func findExisting(email string) bool {
	row := config.DB.QueryRow("SELECT * FROM usuario WHERE email = ?", email)
	if row != nil {
		return true
	}
	return false
}

// AlreadyLoggedIn : checa en la bd si existe una session
func AlreadyLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("blablabla")
	if err != nil {
		return false
	}

	row := config.DB.QueryRow("SELECT * FROM sessions WHERE session_id = ?", cookie.Value)
	if row == nil {
		return false
	}
	return true
}

// InsertIntoSesion : insertar nueva celda a la tabla sesion
func InsertIntoSesion(cValue string, uEmail string) {
	// insertar valores a la tabla sessions
	_, err := config.DB.Exec("INSERT INTO sessions(session_id, user_id) VALUES (?, ?)", cValue, uEmail)
	if err != nil {
		fmt.Println("Error en insertar valores a la tabla sessions...")
		log.Fatalln(err)
	}
}

// DeleteSesion : borrar sesion dada
func DeleteSesion(cValue string) {
	_, err := config.DB.Exec("DELETE FROM sessions WHERE session_id = ?", cValue)
	if err != nil {
		fmt.Println("Error en borrar session de sessions...")
		log.Fatalln(err)
	}
}
