package controllers

import (
	"encoding/json"
	"errors"
	"go-todo-app/config"
	 "go-todo-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.user
	//retrieve the value stored in the "decryptedText" key in the context. If the key is present, it will assign the value stored in the key to the variable decryptedData &
	//set the boolean variable exists to true. If the key is not present, then decryptedData will be set to its zero value and exists will be set to false.
	decryptedData, exists := context.Get("decryptedText")
	if !exists {
		context.AbortWithError(http.StatusBadRequest, errors.New("decrypted data not found"))
		return
	}
	//convert a decrypted data into a user object. The 1st argument is a byte array of decrypted data, and 2nd  argument is the address of the user object.
	//The function will attempt to unmarshal the data into the user object.
	json.Unmarshal(decryptedData.([]byte), &user)

	//connect to db
	db := config.ConnectToDB()
	defer db.Close()
	//insert user data into table
	_, err := db.Query("insert into users(Name, Username, Email, Password) values(?,?,?,?)", user.Name, user.Username, user.Email, user.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, err)

		// context.Abort() func is used to abort the current context. This will stop the current context from continuing and
		//will immediately return control to the caller. This is usually used when an error occurs or when a task needs to be terminated before it completes.
		context.Abort()
		return
	}
	//return a JSON response with an encrypted message. The message is encrypted using the AESEncrypt function, which takes a string, a byte array (x-key),& a string ( x-iv).
	// The response status is set to "Created" (HTTP Status 201), indicating that the request was successful.
	context.JSON(http.StatusCreated, AESEncrypt("Success........", []byte(context.Request.Header.Get("x-key")), context.Request.Header.Get("x-iv")))
}
