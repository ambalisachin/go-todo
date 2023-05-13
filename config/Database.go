package config

import (
	"database/sql"
	"fmt"
)

// ConnectToDB connects to the database 
func ConnectToDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Connected to DB Successfully....... ")
	return db
}

// NewTable creates new table if the table not exist
func NewTable() {
	db := ConnectToDB()
	defer db.Close()
	//checking for create table for user in db exist or not , if not in db
	//crate a table for user in db
	_, err := db.Query("CREATE TABLE IF NOT EXISTS users(Name varchar(20) UNIQUE NOT NULL, Username varchar(20) NOT NULL, Email varchar(20) NOT NULL, Password varchar(20) NOT NULL)")
	if err != nil {
		fmt.Println(err)
	}
	//checking for create table for todo in db exist or not , if not in db
	//creat table for todo in db
	_, e := db.Query("CREATE TABLE IF NOT EXISTS todo(ID int(5) UNIQUE NOT NULL AUTO_INCREMENT, Title varchar(20) NOT NULL, Description varchar(20))")
	if e != nil {
		fmt.Println(e)
	}
}
