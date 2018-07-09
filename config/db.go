package config

import (
	"database/sql"
	"fmt"
	//libreria para mysql
	_ "github.com/go-sql-driver/mysql"
)

// DB : database variable
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:password@/login")
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You are connected to your database...")

}
