package mysql

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() {
	err := InitSnowFlack("2024-01-29", 1)
	if err != nil {
		log.Fatalln("雪花算法初始化失败")
		return
	}
	dsn := "root:1104850836@tcp(127.0.0.1:3306)/BookWeb?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		//链接失败直接退出程序
		log.Fatalln(err)
	}
	DB = db
	DB.SetMaxIdleConns(8)
	DB.SetMaxOpenConns(10)
	err = DB.Ping()
	//defer DB.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
