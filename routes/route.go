package routes

import (
	"go-todo-app/controllers"
	middleware "go-todo-app/middlewares"

	"github.com/gin-gonic/gin"
)

//var encryptedString string

func SetupRouter() *gin.Engine {
	//	creates a router as 'r' & sets it to use the gin framework's default settings.
	//This allows the router to use all of the default routes and middleware functions that are available in the gin framework.
	r := gin.Default()
	//	creates a new router group named 'v1' which is associated with the URL path prefix '/v1'. This router group can be used to handle routes specific to the '/v1' prefix.
	v1 := r.Group("/v1")
	//creates a route group  called "v1". The route group is called "/add" and can contain routes that are related to adding something.
	add := v1.Group("/add")
	//using a middleware called DecryptRequest().
	//This middleware is responsible for decrypting requests that are sent to the server.
	// It will decrypt any encrypted requests that are made to the server, allowing the server to access any data that is encrypted.
	// The middleware will also ensure that all requests are correctly authenticated and authorized, preventing unauthorized access to sensitive data.
	add.Use(middleware.DecryptRequest())
	{
		add.POST("todo", Controllers.CreateATodo)
		add.POST("/user/register", Controllers.RegisterUser)
		add.PUT("todo/:id", Controllers.UpdateATodo)
		v1.GET("todo", Controllers.GetTodos)
		v1.GET("todo/:id", Controllers.GetATodo)
		v1.DELETE("todo/:id", Controllers.DeleteATodo)
		v1.POST("/token", Controllers.Login)

		secured := v1.Group("/secured")
		{
			secured.GET("/ping", Controllers.Ping)
		}
	}
	//creates a route group called "encrypt" with the path of "/data". Any routes within the encrypt route group will have the path prefix of "/data".
	encrypt := r.Group("/data")
	{
		encrypt.POST("encrypt", Controllers.EncryptDataHandler)
		encrypt.POST("decrypt", Controllers.DecryptDataHandler)
	}
	return r
}
