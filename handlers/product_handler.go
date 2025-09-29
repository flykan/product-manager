package handlers

import (
	"github.com/flykan/product-manager/database"
	"github.com/flykan/product-manager/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts 获取所有商品（支持分页和搜索）
func GetProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search[value]")

	offset := (page - 1) * limit

	// 构建查询
	query := database.DB.Model(&models.Product{})

	// 搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR category LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	query.Count(&total)

	// 获取数据
	err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":            products,
		"recordsTotal":    total,
		"recordsFiltered": total,
		"draw":            c.Query("draw"),
	})
}

// GetProduct 获取单个商品
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	err := database.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct 创建商品
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	// 验证必填字段
	if product.Name == "" || product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name and price are required",
		})
		return
	}

	err := database.DB.Create(&product).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct 更新商品
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	err := database.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	// 更新字段
	if updateData.Name != "" {
		product.Name = updateData.Name
	}
	if updateData.Description != "" {
		product.Description = updateData.Description
	}
	if updateData.Price > 0 {
		product.Price = updateData.Price
	}
	if updateData.Stock >= 0 {
		product.Stock = updateData.Stock
	}
	if updateData.Category != "" {
		product.Category = updateData.Category
	}

	err = database.DB.Save(&product).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct 删除商品
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	err := database.DB.First(&product, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	err = database.DB.Delete(&product).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
