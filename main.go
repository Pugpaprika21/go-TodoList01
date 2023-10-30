package main

import (
	"go-TodoList/db"
	"go-TodoList/helper"
	"log"
	"os"
	"text/template"

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
		if err := godotenv.Load(); err != nil { // dev
			log.Fatal("Error Loading .env file")
		}
	}

	db.ConnectDB()
	db.Migrate()

	store := cookie.NewStore([]byte("secret"))

	router := gin.Default()
	router.Use(sessions.Sessions("mysession", store))

	router.SetFuncMap(template.FuncMap{
		"Rows": helper.NumRows,
		"DMY":  helper.GetDMY,
	})

	router.LoadHTMLGlob("view/**/*")
	router.Static("/assets", "./assets")

	serveRoutes(router)

	port := os.Getenv("PORT")
	if port != "" {
		port = "5000"
	}

	router.Run(":" + port)
}
