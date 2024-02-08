package main

import (
	"github.com/gin-gonic/gin"
	auth_controllers "github.com/go-redis-caching-poc/controllers/auth-controllers"
	common_controllers "github.com/go-redis-caching-poc/controllers/common-controllers"
)

func main () {
	server := gin.New();
	
	server.GET("/home", common_controllers.PrintHelloWorld)
	server.POST("/sign-up", auth_controllers.SignUp)
	server.POST("/login", auth_controllers.LogIn)
	server.POST("/refresh-token", auth_controllers.RefreshToken)

	server.Run(":8085")
}