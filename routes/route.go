package routes

import (
	"go-todo-app/controllers"
	"go-todo-app/middlewares"

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
		add.POST("todo", controllers.CreateATodo)
		add.POST("/user/register", controllers.RegisterUser)
		add.PUT("todo/:id", controllers.UpdateATodo)
		v1.GET("todo", controllers.GetTodos)
		v1.GET("todo/:id", controllers.GetATodo)
		v1.DELETE("todo/:id", controllers.DeleteATodo)
		v1.POST("/token", controllers.Login)

		secured := v1.Group("/secured")
		{
			secured.GET("/ping", controllers.Ping)
		}
	}
	//creates a route group called "encrypt" with the path of "/data". Any routes within the encrypt route group will have the path prefix of "/data".
	encrypt := r.Group("/data")
	{
		encrypt.POST("encrypt", controllers.EncryptDataHandler)
		encrypt.POST("decrypt", controllers.DecryptDataHandler)
	}
	return r
}
