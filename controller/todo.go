package controller

import (
	"go-TodoList/db"
	"go-TodoList/model"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Todo struct{}

var appUrl = os.Getenv("APP_URL") + os.Getenv("APP_NAME")

func (t *Todo) Index(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set("username", "alex")
	session.Save()

	ctx.HTML(http.StatusOK, "todo.html", gin.H{
		"url":   appUrl + "/create",
		"title": "Todo-list",
	})
}

func (t *Todo) Create(ctx *gin.Context) {
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	session := sessions.Default(ctx)
	username := session.Get("username").(string)

	user := model.User{
		Username: username,
		Password: "1234",
		Status:   1,
	}
	db.Conn.Create(&user)

	todo := model.Todo{
		Title:       title,
		Description: description,
		UserID:      user.ID,
	}

	db.Conn.Create(&todo)

	ctx.Redirect(http.StatusCreated, appUrl+"/index")
}
