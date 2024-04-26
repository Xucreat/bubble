package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Define a global variable
var (
	DB *gorm.DB
)

// Todo model
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// connect database
func initMySQL() (err error) {
	dsn := "bx:123456@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn) // return a database object and an error
	if err != nil {
		return
	}
	return DB.DB().Ping() // connect->nil; disconnect->err
}

func main() {
	// create db
	// sql: CREATE DATABASE bubble
	// connect database
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close() // The program exits and closes the database connection

	// model binding
	DB.AutoMigrate(&Todo{}) // todos

	r := gin.Default()

	// Tell where the static files referenced by the template file are
	r.Static("/static", "static")

	// Tell gin where to find templates
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// v1
	v1Group := r.Group("v1")
	{
		// to do list
		// add
		v1Group.POST("/todo", func(c *gin.Context) {
			// 1. 从请求中把数据拿出来
			var todo Todo
			c.BindJSON(&todo) // 从请求中读取 JSON 数据，并将其绑定到 todo 变量上。这意味着请求中的 JSON 数据将被解析，并将其中的字段值分配给 todo 结构体的对应字段。

			// 2. 存入数据库
			err = DB.Create(&todo).Error

			// 3. 返回响应
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}

		})

		// check all
		v1Group.GET("/todo", func(c *gin.Context) {
			// check all datas in tables
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		// check one
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})

		// update one
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "volid id"})
			}
			var todo Todo
			if err = DB.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			}
			// ↑检验前端传来的数据
			c.BindJSON(&todo) // 将前端请求的数据绑定到todo
			// ↓更新数据库中的数据
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}

		})

		// delete one
		v1Group.DELETE("todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "volid id"})
			}
			if err = DB.Where("id = ?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "deleted"})
			}
		})
	}

	r.Run()
}
