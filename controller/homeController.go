package controller

import (
	"Delfin/config"
	"net/http"
)

// Home : Landing Page
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "home.html", nil)
}
