package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-todo-app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req

	// connect to the test database
	//db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/testdb")
	//db := ConnectToDB()
	//config.NewTable()
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer db.Close()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	defer db.Close()

	// create a new todo record in the database
	_, err = db.Exec("INSERT INTO todo (title, description) VALUES (?, ?)", "Test Todo", "Test Description")
	if err != nil {
		t.Fatal(err)
	}

	// call the GetTodos function
	GetTodos(ctx)

	// check the response
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	// verify that the response contains the expected data
	var todos []models.Todo
	err = json.Unmarshal([]byte(rr.Body.String()), &todos)
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 1 {
		t.Errorf("expected 1 todo; got %v", len(todos))
	}
	if todos[0].Title != "Test Todo" {
		t.Errorf("expected todo title 'Test Todo'; got %v", todos[0].Title)
	}
	if todos[0].Description != "Test Description" {
		t.Errorf("expected todo description 'Test Description'; got %v", todos[0].Description)
	}
}

func TestCreateATodo(t *testing.T) {
	
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer db.Close()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	defer db.Close()

	// Create a test request with sample encrypted data
	encryptedData := []byte(`{"ID": 1, "Title": "Test Title", "Description": "Test Description"}`)
	req, err := http.NewRequest(http.MethodPost, "/todo", nil)
	req.Header.Set("x-key", "sample-key")
	req.Header.Set("x-iv", "sample-iv")
	req.SetBasicAuth("username", "password")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Set up a test context with the encrypted data
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set("decryptedText", encryptedData)
	ctx.Request = req

	// Call the CreateATodo function with the test context
	CreateATodo(ctx)

	// Assert that the response is a successful HTTP status and contains the expected message
	assert.Equal(t, http.StatusCreated, w.Code)
	expected := AESEncrypt("Todo created Successfully.....", []byte("sample-key"), "sample-iv")
	actual := w.Body.String()
	assert.Equal(t, expected, actual)

}
func TestGetATodo(t *testing.T) {
	// Create a new Gin router instance
	r := gin.Default()

	// Add a GET route for testing the GetATodo function
	r.GET("/todo/:id", GetATodo)
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer db.Close()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	defer db.Close()

	// Create a test request with a sample todo ID
	req, err := http.NewRequest(http.MethodGet, "/todo/1", nil)
	req.Header.Set("x-key", "sample-key")
	req.Header.Set("x-iv", "sample-iv")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Call the GetATodo function with the test context
	r.ServeHTTP(w, req)

	// Assert that the response is a successful HTTP status and contains the expected todo data
	assert.Equal(t, http.StatusOK, w.Code)
	expectedData := `{"ID":1,"Title":"Test Title","Description":"Test Description"}`
	expected := AESEncrypt(expectedData, []byte("sample-key"), "sample-iv")
	actual := w.Body.String()
	assert.Equal(t, expected, actual)

}
