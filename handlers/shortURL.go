package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Title string
}

func HomeRouteHandler(g *gin.Context) {
	var requestBody RequestBody = RequestBody{}

	// if err := g.BindJSON(&requestBody); err != nil {
	// 	g.AbortWithStatus(http.StatusBadRequest)
	// }

	readerBytes, err := io.ReadAll(g.Request.Body)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Errorf(err.Error()),
		})
		return
	}

	if err := json.Unmarshal(readerBytes, &requestBody); err != nil {
		fmt.Println(err)
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Errorf(err.Error()),
		})
		return
	}

	logger := log.Default()

	logger.Println(requestBody)

	fmt.Println(requestBody.Title)

	g.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully hit with title: %s", requestBody.Title),
	})
}
