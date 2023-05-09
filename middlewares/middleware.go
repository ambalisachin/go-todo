package middleware

import (
	"encoding/base64"
	"fmt"
	"go-todo-app/controllers"

	"github.com/gin-gonic/gin"
)

// DecryptRequest func returns a gin.HandlerFunc. This handler func can be used to decrypt the request body.
// This function when it receives a request and the provided handler function will be executed to decrypt the request body and make it available for further processing.
func DecryptRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//creates a variable called "requestBody" which is a map of strings. the keys & the values are strings.
		//This variable can be used to store data which is sent in an HTTP request body.
		var requestBody map[string]string
		//Bind the request body of a request to a struct, named requestBody. The ShouldBindJSON() method  is used to read the request body and bind it to the struct.
		//The ctx is a context object that stores the request and response objects.
		ctx.ShouldBindJSON(&requestBody)
		fmt.Println(requestBody, requestBody["Encrypted-data"])
		//it takes a string representation of data that has been encrypted and decodes it into a byte array format.
		//The requestBody variable is assumed to be a map containing the key "Encrypted-data" with a value that is a string.
		// The DecodeString function is part of the Golang base64 package and it takes in a string and returns a byte array and an error.

		encryptedData, _ := base64.StdEncoding.DecodeString(requestBody["Encrypted-data"])
		//decrypt data that has been encrypted using AES encryption. The data which is encrypted is passed as 1st argument to the Controllers.
		//AESDecrypt func & 2nd argument is a byte slice containing the encryption key, which has been retrieved from the request header with the "x-key" key.
		//3rd argument is the IV, which has been retrieved from the request header with the "x-iv" key. The decrypted data is stored in the decryptedText variable.
		decryptedText := controllers.AESDecrypt(encryptedData, []byte(ctx.Request.Header.Get("x-key")), ctx.Request.Header.Get("x-iv"))
		fmt.Println("\n decrypted data:", string(decryptedText))

		//Setting a value to the context variable "decryptedText" with the value of the variable "decryptedText".
		//This allows the value of decryptedText to be available throughout the program, making it easy to access and use for any other functions.
		ctx.Set("decryptedText", decryptedText)

		//it allows for the execution of the next handler in the chain. It is used to implement middleware patterns such as authentication and logging. It is typically used in a handler function and will be called after all the necessary work has been completed.
		//This gives the next handler in the chain the chance to execute and return a response.
		ctx.Next()
	}
}
