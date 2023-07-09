package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleBasicRoute(c *gin.Context) {
	fmt.Println("this is the basic route handler")
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully handled route",
	})
}

func main() {
	const PORT = 8080
	fmt.Println("This is a new go module")
	router := gin.Default()
	router.GET("/", handleBasicRoute)

	router.Run(fmt.Sprintf(":%d", PORT))
}
