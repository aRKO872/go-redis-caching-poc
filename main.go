package main

import (
	"context"

	"github.com/gin-gonic/gin"
	auth_controllers "github.com/go-redis-caching-poc/controllers/auth-controllers"
	common_controllers "github.com/go-redis-caching-poc/controllers/common-controllers"
	db_controllers "github.com/go-redis-caching-poc/controllers/db-controllers"
	redis_controllers "github.com/go-redis-caching-poc/controllers/redis-controllers"
)

func init() {
	db_controllers.DatabaseConnection()
	redis_controllers.InitializeRedis(context.Background())
}

func main () {
	server := gin.Default();
	
	server.GET("/home", common_controllers.PrintHelloWorld)

	server.POST("/sign-up", auth_controllers.SignUp)
	server.POST("/login", auth_controllers.LogIn)
	server.POST("/refresh-token", auth_controllers.RefreshToken)

	server.POST("/update-details", common_controllers.UpdateDetails)
	server.GET("/get-all-user-detail", common_controllers.GetAllUserDetails)

	server.Run(":8085")
}