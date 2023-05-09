package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// logging and authorization
var jwtKey = []byte("secret_key")

// Store & access user information
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// This code creates a struct named Claims which contains a field named Username of type string and a field of type jwt.StandardClaims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	// decode JSON data sent in an HTTP request.
	// var "credentials" is of type "Credentials".json.NewDecoder() is used to create a new decoder object which will be used to decode the JSON data sent in the request.
	//The decoded data is then stored in the "credentials" variable. If there is an error while decoding the data, the http status code "400" is sent to the client and the code execution is stopped.
	var credentials Credentials
	err := json.NewDecoder(c.Request.Body).Decode(&credentials)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}
