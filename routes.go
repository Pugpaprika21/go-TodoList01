package main

import (
	"go-TodoList/controller"

	"github.com/gin-gonic/gin"
)

func serveRoutes(router *gin.Engine) {
	todoController := controller.Todo{}
	todoGroup := router.Group("/todo")

	todoGroup.GET("/login", todoController.Login)
	todoGroup.POST("/checkLogin", todoController.CheckLogin)
	todoGroup.GET("/index", todoController.Index)
	todoGroup.POST("/create", todoController.Create)
	todoGroup.GET("/logout", todoController.Logout)
}
