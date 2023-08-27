package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anarchymonkey/url-shortner/database"
	"github.com/anarchymonkey/url-shortner/globals"
	"github.com/anarchymonkey/url-shortner/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Env struct {
	DB  *gorm.DB
	env string
}

func getDatabaseConfig(dbName string) database.DatabaseConfig {
	return database.DatabaseConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: os.Getenv("PG_PASS_V1"),
		Dbname:   dbName,
		Port:     5432,
	}

}

func handlerWrapper(fn func(*gin.Context, *gorm.DB), env *Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fn(ctx, env.DB)
	}
}

func main() {
	router := gin.Default()

	router.Use(corsMiddleWare())

	db, err := database.Connect(getDatabaseConfig("url_shortner"))

	if err != nil {
		fmt.Errorf("Error while initializing database connection")
		os.Exit(1)
	}

	var env *Env = &Env{
		env: "local",
		DB:  db,
	}

	router.POST("/generate-short-url", handlerWrapper(handlers.GenerateShortURL, env))

	router.Run(globals.SERVER_PORT)
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
