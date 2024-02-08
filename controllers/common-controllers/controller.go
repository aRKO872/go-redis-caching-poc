package common_controllers

import "github.com/gin-gonic/gin"

func PrintHelloWorld(c *gin.Context) {
	c.JSON (200, gin.H{
		"status": true,
		"message": "Hello Word!",
	})
}