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

func (t *Todo) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"url":   appUrl + "/create",
		"title": "Todo-list",
	})
}

func (t *Todo) CheckLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	var user model.User
	if err := db.Conn.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		user.Username = username
		user.Password = password
		user.Status = 1
		db.Conn.Create(&user)
		ctx.Redirect(http.StatusSeeOther, "/todo/login")
		return
	}
	session := sessions.Default(ctx)
	session.Set("userId", user.ID)
	session.Set("username", user.Username)
	session.Set("password", user.Password)
	session.Set("status", user.Status)
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/todo/index")
}

func (t *Todo) Index(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get("userId")
	username := session.Get("username")
	password := session.Get("password")

	if userId == nil {
		ctx.Redirect(http.StatusSeeOther, "/todo/login")
		return
	}

	var todos []model.Todo
	db.Conn.Where("user_id = ?", userId).Find(&todos)

	ctx.HTML(http.StatusOK, "todo.html", gin.H{
		"user": gin.H{
			"userId":   userId,
			"username": username,
			"password": password,
		},
		"todos": todos,
		"title": "Todo-list",
	})
}

func (t *Todo) Create(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get("userId")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	if userId == nil {
		ctx.Redirect(http.StatusSeeOther, "/todo/login")
		return
	}

	todo := model.Todo{
		Title:       title,
		Description: description,
		UserID:      userId.(uint),
	}

	db.Conn.Create(&todo)
	ctx.Redirect(http.StatusSeeOther, "/todo/index")
}

func (t *Todo) LogOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/todo/login")
}
