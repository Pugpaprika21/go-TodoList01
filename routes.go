package main

import (
	"go-TodoList/controller"

	"github.com/gin-gonic/gin"
)

func serveRoutes(router *gin.Engine) {
	todoGroup := router.Group("/todo")
	todoController := controller.Todo{}
	todoGroup.GET("/login", todoController.Login)
	todoGroup.POST("/checkLogin", todoController.CheckLogin)
	todoGroup.GET("/logout", todoController.Logout)
	todoGroup.GET("/index", todoController.Index)
	todoGroup.POST("/create", todoController.Create)
	todoGroup.PATCH("/update/:id", todoController.Update)
	todoGroup.DELETE("/delete/:id", todoController.Delete)
}
