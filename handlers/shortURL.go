package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	LongURL string
}

func parseRequestBody[T comparable](g *gin.Context, request T) (T, error) {
	readerBytes, err := io.ReadAll(g.Request.Body)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return request, fmt.Errorf(fmt.Sprintf("Error occoured %v", err))
	}

	if err := json.Unmarshal(readerBytes, &request); err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return request, fmt.Errorf(fmt.Sprintf("Error occoured %v", err))
	}

	return request, nil
}

func GenerateShortURL(g *gin.Context) {
	var requestBody RequestBody = RequestBody{}
	requestBody, err := parseRequestBody(g, requestBody)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Error parsing the request body %v", err),
		})
		return
	}

	fmt.Println(requestBody.LongURL)

	g.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully hit with title: %s", requestBody.LongURL),
	})
}
