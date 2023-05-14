package config

import (
	"database/sql"
	"fmt"
	"time"
)

// Credentials struct can be used to store credentials in a single data type.
type Credentials struct {
	Username string
	Password string
	Server   string
	Dbname   string
}

var Database = Credentials{
	Username: "root",
	Password: "password",
	Server:   "tcp(localhost:3306)",
	Dbname:   "sachindb",
}

// ConnectToDB connects to the database
func (m Credentials) ConnectToDB() *sql.DB {
	dataSourceName := m.Username + ":" + m.Password + "@" + m.Server + "/" + m.Dbname
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour * 1)
	fmt.Println("Connected to DB Successfully....... ")
	return db
}

// NewTable creates new table if the table not exist
func NewTable() {
	db := Database.ConnectToDB()
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
