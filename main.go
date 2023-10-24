package main

import (
	"go-TodoList/db"
	"log"
	"os"

	//"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error Loading .env file")
		}
	}

	db.ConnectDB()
	db.Migrate()

	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("view/**/*")
	r.Static("/assets", "./assets")

	serveRoutes(r)

	port := os.Getenv("PORT")
	if port != "" {
		port = "5000"
	}

	r.Run(":" + port)
}
