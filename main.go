package main

import (
	"github.com/flykan/product-manager/database"
	"github.com/flykan/product-manager/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 创建Gin路由
	router := gin.Default()

	// 静态文件服务
	router.Static("/static", "./static")
	router.LoadHTMLGlob("static/*.html")

	// 路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// API路由
	api := router.Group("/api")
	{
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.POST("/products", handlers.CreateProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)
	}

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
