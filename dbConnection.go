package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func testMysqlConnection(db *sql.DB) {
	var err error
	db, err = sql.Open("mysql", "root:sakila@tcp(127.0.0.1:3306)/sakila")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Mysql connection successful")
	defer db.Close()
}
