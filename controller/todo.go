package controller

import (
	"go-TodoList/db"
	"go-TodoList/model"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Todo struct{}

func (t *Todo) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Todo-list",
	})
}

func (t *Todo) CheckLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username == "" && password == "" {
		ctx.Redirect(http.StatusSeeOther, "/todo/login")
		return
	}

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

	currentTime := time.Now()
	dmyFormat := currentTime.Format("2006-01-02")

	if userId == nil {
		ctx.Redirect(http.StatusSeeOther, "/todo/login")
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var perPage = 10
	var bookCount int64
	db.Conn.Model(&model.Todo{}).Where("created_at IS NOT NULL AND user_id = ?", userId).Count(&bookCount)

	pageCount := int(math.Ceil(float64(bookCount) / float64(perPage)))

	if page < 1 || page > pageCount {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	offset := (page - 1) * perPage

	var todos []model.Todo
	db.Conn.Where("created_at IS NOT NULL AND user_id = ?", userId).
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&todos)

	prevPage := page - 1
	nextPage := page + 1

	var pages []int
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, i)
	}

	ctx.HTML(http.StatusOK, "todo.html", gin.H{
		"user": gin.H{
			"userId":   userId,
			"username": username,
			"password": password,
		},
		"nowDMY":       dmyFormat,
		"todos":        todos,
		"currentPage":  page,
		"totalCount":   int(bookCount),
		"itemsPerPage": perPage,
		"pageCount":    pageCount,
		"prevPage":     prevPage,
		"nextPage":     nextPage,
		"pages":        pages,
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
	ctx.Redirect(http.StatusSeeOther, "/todo/index?msg=created")
}

func (t *Todo) Edit(ctx *gin.Context) {
	todoId := ctx.Param("id")
	session := sessions.Default(ctx)
	userId := session.Get("userId")

	ctx.String(http.StatusAccepted, todoId+""+userId.(string))
}

func (t *Todo) Delete(ctx *gin.Context) {
	todoId := ctx.Param("id")
	db.Conn.Delete(&model.Todo{}, todoId)
	ctx.JSON(http.StatusOK, gin.H{
		"todoId": todoId,
	})
}

func (t *Todo) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/todo/login")
}
