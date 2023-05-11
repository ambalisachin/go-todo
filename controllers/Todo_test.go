package controllers

import (
	"bytes"
	"database/sql"
	
	"fmt"
	//"go-todo-app/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTodos(t *testing.T) {
	// create a new gin context
	//r := gin.Default()
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
	// Create a new Gin router instance
	//r := gin.Default()
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

func TestUpdateATodo(t *testing.T) {
	// initialize Gin router
	//r := gin.Default()
	// set up test database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// insert a test record into the database
	_, err = db.Exec("insert into todo (ID, Title, Description) values (?, ?, ?)", 1, "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("failed to insert test record: %v", err)
	}

	// create a mock HTTP request with a sample encrypted JSON payload
	req, err := http.NewRequest("PUT", "/todo/1", bytes.NewBufferString(`{"Title":"Updated Title","Description":"Updated Description"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("x-key", "sample-key")
	req.Header.Set("x-iv", "sample-iv")
	// set up mock HTTP response recorder
	resp := httptest.NewRecorder()

	// set up mock Gin context
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	// simulate adding the decrypted data to the Gin context
	c.Set("decryptedText", []byte(`{"Title":"Updated Title","Description":"Updated Description"}`))

	// call the handler function
	UpdateATodo(c)

	// check the response status code
	if resp.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Code)
	}

	// check the response body
	expected := `{"data":"VXBkYXRlZCBTdWNjZXNzZnVsLi4uLi4uLi4sLi4uLi4uLi4sLi4uLi4uLi4uLi4sLi4uLi4uLi4="}`
	if resp.Body.String() != expected {
		t.Errorf("expected body %v; got %v", expected, resp.Body.String())
	}

	// check that the record was updated in the database
	var title, desc string
	err = db.QueryRow("select Title, Description from todo where ID = ?", 1).Scan(&title, &desc)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if title != "Updated Title" || desc != "Updated Description" {
		t.Errorf("expected record to be updated; got Title=%v, Description=%v", title, desc)
	}
}
func TestDeleteATodo(t *testing.T) {
	// initialize Gin router
	//router := gin.Default()
	// set up test database
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/sachindb")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// insert a test record into the database
	_, err = db.Exec("insert into todo (ID, Title, Description) values (?, ?, ?)", 1, "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("failed to insert test record: %v", err)
	}

	// create a mock HTTP request
	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("x-key", "sample-key")
	req.Header.Set("x-iv", "sample-iv")
	// set up mock HTTP response recorder
	resp := httptest.NewRecorder()

	// set up mock Gin context
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	// call the handler function
	DeleteATodo(c)

	// check the response status code
	if resp.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Code)
	}

	// check the response body
	expected := `{"data":"UmVjb3JkIGRlbGV0ZWQgU3VjY2Vzc2Z1bC4uLi4uLi4uLi4sLi4uLi4uLi4uLi4uLi4uLi4sLi4uLi4uLi4="}`
	if resp.Body.String() != expected {
		t.Errorf("expected body %v; got %v", expected, resp.Body.String())
	}

	// check that the record was deleted from the database
	var count int
	err = db.QueryRow("select count(*) from todo where ID = ?", 1).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if count != 0 {
		t.Errorf("expected record to be deleted; got count=%v", count)
	}
}
