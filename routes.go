package main

import (
	"go-TodoList/controller"

	"github.com/gin-gonic/gin"
)

func serveRoutes(r *gin.Engine) {
	todoController := controller.Todo{}
	todoGroup := r.Group("/todo")

	todoGroup.GET("/login", todoController.Login)
	todoGroup.POST("/checkLogin", todoController.CheckLogin)
	todoGroup.GET("/index", todoController.Index)
	todoGroup.POST("/create", todoController.Create)
	todoGroup.GET("/LoginOut", todoController.LoginOut)
}
