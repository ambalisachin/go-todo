package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go-todo-app/config"
	"go-todo-app/models"

	"github.com/gin-gonic/gin"
)

// GetTodos is used to define a function in Golang which is used to get all the todos from a database.
func GetTodos(c *gin.Context) {
	//GetTodos function in Golang that handles a GET request for a list of todos.
	var todos []models.Todo
	// db := Config.ConnectToDB()
	db := config.Database.ConnectToDB()
	defer db.Close()
	//Query a database table called "todo".
	//db.Query() func to query the table and stores the results in a variable called "row".
	//checks if an error occurred while querying the table. If an error occurred,
	//the code will print the error message to the console and then return.
	row, err := db.Query("SELECT * FROM todo")
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	//Iterates over a collection of rows from a SQL query and stores each row into the "todo" variable which is of type Models.Todo.
	//It does this by scanning each row and assigning the values to the ID, Title, and Description fields of the todo variable.
	// If an error is encountered, the error is printed to the writer.
	for row.Next() {
		var todo models.Todo
		if err := row.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			fmt.Fprint(c.Writer, err)
			return
		}
		//Adds a "todo" item to the list of "todos".
		//aappend func takes 2 arguments: the list of existing todos and the new todo item that is to be added to the list.
		//Func then adds the new todo item to the end of the existing list and returns the new list.
		todos = append(todos, todo)
	}
	data, _ := json.Marshal(todos)
	fmt.Println(string(data))
	//Send an HTTP response with an array of "todos" as the body of the response,and
	//a status code of 200 (OK). func c.JSON() is used to respond with JSON and the "todos" is the data which will be sent in the response body.
	//The HTTP status code of 200 indicates that the request was successful.
	c.JSON(http.StatusOK, AESEncrypt(string(data), []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv")))
}

// CreateATodo function creates a new todo item .
func CreateATodo(c *gin.Context) {
	//The var "todo" is of type Models.Todo, which is a type defined in the Models package.
	//This variable can be used to store data related to a Todo type, such as its title, description, and completion status.
	var todo models.Todo
	decryptedData, exists := c.Get("decryptedText")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errors.New("decrypted data not found"))
		return
	}
	json.Unmarshal(decryptedData.([]byte), &todo)
	db := config.Database.ConnectToDB()
	defer db.Close()
	//Trying to add a new record to a database table called "todo".
	//Query() func from the db package to execute an SQL query. The query is an INSERT statement that adds a new record to the todo table.
	//The values for the record are provided as parameters, including the todo ID, title, and description.
	//If the query is unsuccessful, an error is returned and the code returns a Bad Request response with the error.
	_, err := db.Query("insert into todo(Title, Description) values(?,?)", todo.Title, todo.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, AESEncrypt("Todo created Successfully.....", []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv")))
}

// GetATodo is a function that retrieves a to-do item from a database
func GetATodo(c *gin.Context) {
	//assign the value of the "id" parameter from the "c" object to a var called "id"."c" object is assumed to be an instance of a type that provides access to the "Params" object.
	//The "Params" object is assumed to have a method called "ByName" which takes a parameter and returns the value of the corresponding parameter from the "c" object.

	id := c.Params.ByName("id")
	var todo models.Todo
	db := config.Database.ConnectToDB()
	defer db.Close()
	row, err := db.Query("SELECT * FROM todo where ID=?", id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	for row.Next() {
		if err := row.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			fmt.Fprint(c.Writer, err)
			return
		}
	}
	data, _ := json.Marshal(todo)
	fmt.Println(data)
	c.JSON(http.StatusOK, AESEncrypt(string(data), []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv")))
}

// UpdateATodo updates an existing todo item.
func UpdateATodo(c *gin.Context) {
	id := c.Params.ByName("id")
	var todo models.Todo
	decryptedData, exists := c.Get("decryptedText")
	if !exists {
		c.AbortWithError(http.StatusBadRequest, errors.New("decrypted data not found"))
		return
	}
	json.Unmarshal(decryptedData.([]byte), &todo)
	db := config.Database.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("update todo set Title=?, Description=? where ID=?", todo.Title, todo.Description, id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	c.JSON(http.StatusOK, AESEncrypt("Updated Successfully.......", []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv")))
}

// DeleteATodo function deletes a todo item from the todos table in a database using the given ID
func DeleteATodo(c *gin.Context) {
	id := c.Params.ByName("id")
	db := config.Database.ConnectToDB()
	defer db.Close()
	_, err := db.Exec("DELETE from todo where ID=?", id)
	if err != nil {
		fmt.Fprint(c.Writer, err)
		return
	}
	c.JSON(http.StatusOK, AESEncrypt("Record deleted Succesfully.......", []byte(c.Request.Header.Get("x-key")), c.Request.Header.Get("x-iv")))
}
