package main

import (
	"go-todo-app/config"
	"go-todo-app/routes"

	_ "github.com/go-sql-driver/mysql"
)

//var err error

func main() {
	//this func creates a new table in the configuration. This table can be used to store data related to
	//the configuration, such as settings, values, etc. It can also be used to access and modify the configuration data.

	config.NewTable()

	//The SetupRouter() func is used to create a new router and assign it to the variable r.
	//This router can then be used to handle requests, create routes, and other functionality.

	r := routes.SetupRouter()

	r.Run("localhost:6000")
}
