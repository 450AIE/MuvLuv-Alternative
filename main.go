package main

import (
	"Web/dao/mysql"
	"Web/router"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	mysql.InitDB()
	router.InitRouter()
	defer mysql.DB.Close()
}
