package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong2",
	})
}



func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello1",
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", pingHandler)

	router.GET("/Hello", helloHandler)

	router.Run(":8080")
}

