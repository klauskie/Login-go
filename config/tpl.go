package config

import "html/template"

//TPL Variable que almacena todos los html
var TPL *template.Template

func init() {
	TPL = template.Must(template.ParseGlob("templates/*.html"))
}
