package controller

import (
	"fmt"
	"go-TodoList/db"
	"go-TodoList/dto"
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
	if err != nil || page < 1 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var perPage = 12
	var todoCount int64
	db.Conn.Model(&model.Todo{}).Where("created_at IS NOT NULL AND user_id = ?", userId).Count(&todoCount)

	pageCount := int(math.Ceil(float64(todoCount) / float64(perPage)))

	if page > pageCount {
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

	if nextPage > pageCount {
		nextPage = 0
	}

	var pages []int
	for i := 1; i <= pageCount; i++ {
		pages = append(pages, i)
	}

	userData := gin.H{
		"userId":   userId,
		"username": username,
		"password": password,
	}

	fileTimestamp := time.Now().Format("20060102150405")

	assetsURL := make(map[string]string)
	assetsURL["css"] = fmt.Sprintf("/assets/css/main.css?v=%s", fileTimestamp)
	assetsURL["js"] = fmt.Sprintf("/assets/js/main.js?v=%s", fileTimestamp)

	ctx.HTML(http.StatusOK, "todo.html", gin.H{
		"user":         userData,
		"nowDMY":       dmyFormat,
		"todos":        todos,
		"currentPage":  page,
		"totalCount":   int(todoCount),
		"itemsPerPage": perPage,
		"pageCount":    pageCount,
		"prevPage":     prevPage,
		"nextPage":     nextPage,
		"pages":        pages,
		"assetsURL":    assetsURL,
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
		Active:      "N",
		UserID:      userId.(uint),
	}

	db.Conn.Create(&todo)
	ctx.Redirect(http.StatusSeeOther, "/todo/index?msg=created")
}

func (t *Todo) Update(ctx *gin.Context) {
	todoId := ctx.Param("id")
	if db.Conn.Model(&model.Todo{}).Where("id = ?", todoId).Update("active", "Y").Error != nil {
		ctx.JSON(http.StatusOK, dto.TodoResponse{
			Status:  200,
			Message: "ปิดงานไม่สำเร็จ",
		})
	}

	ctx.JSON(http.StatusOK, dto.TodoResponse{
		Status:  200,
		Message: "ปิดงานสำเร็จ",
	})
}

func (t *Todo) Delete(ctx *gin.Context) {
	todoId := ctx.Param("id")
	if db.Conn.Delete(&model.Todo{}, todoId).Error != nil {
		ctx.JSON(http.StatusOK, dto.TodoResponse{
			Status:  200,
			Message: "ลบข้อมูลไม่สำเร็จ",
		})
	}

	ctx.JSON(http.StatusOK, dto.TodoResponse{
		Status:  200,
		Message: "ลบข้อมูลสำเร็จ",
	})
}

func (t *Todo) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/todo/login")
}
