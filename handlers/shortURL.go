package handlers

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestBody struct {
	LongURL string
}

var lock sync.Mutex = sync.Mutex{}

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

/*
The functionality should do two things

1. gnerate the short URL
2. dump the short URL along with the long URL and other properties in the DB
*/

func extractBits(hash string, extractBitsCount int) string {
	bits := ""

	for _, b := range hash {
		bits += fmt.Sprintf("%08b", b)
	}

	extractedBits := bits[:extractBitsCount]
	return extractedBits
}

func encodeToBase64(hashToEncode string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(hashToEncode))
}

func generateShortURL(longURL string) string {
	lock.Lock()
	defer lock.Unlock()
	// generate a sha256 hash
	// extract 43 bits
	// encode to base 64

	var hash hash.Hash = crypto.SHA256.New()

	hash.Write([]byte(longURL))
	finalHash := hash.Sum(nil)
	extractedBits := extractBits(string(finalHash), 43)

	fmt.Printf("%x\n %v", finalHash, encodeToBase64(extractedBits))

	return encodeToBase64(extractedBits)[:7]

}
func GenerateShortURL(g *gin.Context, db *gorm.DB) {
	var requestBody RequestBody = RequestBody{}
	requestBody, err := parseRequestBody(g, requestBody)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Error parsing the request body %v", err),
		})
		return
	}

	var expiresAt *time.Time

	shortUrlInsert := ShortenedUrls{
		Longurl:   requestBody.LongURL,
		Shorturl:  generateShortURL(requestBody.LongURL),
		ExpiresAt: expiresAt,
	}

	// creation handled gracefully
	rows := db.Create(&shortUrlInsert)
	fmt.Println("The error is", rows.Error)
	if rows.RowsAffected == 0 {
		shortenedUrls := ShortenedUrls{}
		db.First(&shortenedUrls, "longurl = ?", shortUrlInsert.Longurl)
		g.IndentedJSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("your generated short url is: %s", shortenedUrls.Shorturl),
		})
		return
	}

	fmt.Println(requestBody.LongURL, rows.RowsAffected)

	g.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("your generated short url is: %s", shortUrlInsert.Shorturl),
	})
}

func GetLongURLFromShortURL(g *gin.Context, db *gorm.DB) {
	params, ok := g.Params.Get("shorturl")

	if !ok {
		g.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "query param not found, require short url",
		})
		return
	}

	fmt.Println("the params is", params)
	shortUrlResp := ShortenedUrls{}

	db.First(&shortUrlResp, "shorturl = ?", params)

	fmt.Println(shortUrlResp.Shorturl, shortUrlResp.Longurl, db.RowsAffected, db.Error)
	if db.Error != nil {
		g.IndentedJSON(http.StatusOK, gin.H{
			"message": "long url not found or expired",
		})
		return
	}

	g.Redirect(http.StatusTemporaryRedirect, shortUrlResp.Longurl)
}
