package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Ping Func that takes a pointer to a gin.Context as an argument. The func will perform some action on the gin.Context object, such as setting headers or sending a response.
func Ping(c *gin.Context) {
	//check if there is a token cookie in the context.If there is no token cookie, it returns a status code of 401 (Unauthorized).
	//If there is an error in getting the cookie, it returns a status code of 400 (Bad Request).
	cookie, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// parse a cookie value and setting it to a Claims struct.
	//The Claims struct holds information about the cookie such as the user's name, email address, and other information.
	//The tokenStr is the cookie value that is being parsed and claims is a pointer to a Claims struct which will store the parsed information.
	tokenStr := cookie

	claims := &Claims{}
	//parse a JWT and extract its claims.The tokenStr is a string that contains the JWT.
	// The claims is a struct that will be populated with the token's claims.
	//The jwtKey is a byte array containing the secret key used to sign the JWT.
	// returns the parsed token (tkn) and an error (err) if any.

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	//check if an error has occurred. If an error has occurred,
	// it checks to see if the error is "jwt.ErrSignatureInvalid", which is a specific type of error.
	//If it is, the code writes a status of "Unauthorized" to the response.
	//Otherwise, the code writes a status of "Bad Request" to the response.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//check if the token is valid & if it's not valid, it will return a status code of 401 (unauthorized) to client.
	// If the token is valid, it will return a status code of 200 (OK) and a "pong" message.
	if !tkn.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
