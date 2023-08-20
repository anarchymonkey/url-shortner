package main

import (
	"fmt"
	"net/http"

	"github.com/anarchymonkey/url-shortner/globals"
	"github.com/anarchymonkey/url-shortner/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("This is the URL shortner")

	router := gin.Default()

	router.Use(corsMiddleWare())

	router.POST("/", handlers.HomeRouteHandler)

	router.Run(":8080")
}

func corsMiddleWare() gin.HandlerFunc {

	return func(g *gin.Context) {

		g.Writer.Header().Set("Access-Control-Allow-Origin", globals.ALLOWED_ORIGINS)
		g.Writer.Header().Set("Access-Control-Allow-Methods", globals.ALLOWED_METHODS)
		g.Writer.Header().Set("Access-Control-Allow-Headers", globals.ALLOWED_HEADERS)

		// if there is a options request then follow through to return OK
		if g.Request.Method == http.MethodOptions {
			g.AbortWithStatus(http.StatusOK)
			return
		}
		g.Next()
	}
}
