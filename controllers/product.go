package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
)

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := c.Get("currentUser")
	product.CreatedBy = user.(models.User).ID

	if err := initializers.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product})
}

func GetAllProducts(c *gin.Context) {
	var products []models.Product
	if err := initializers.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func UpdateProducts(c *gin.Context) {
	id := c.Param("id")
	var allProducts models.Product

	if err := c.ShouldBindJSON(&allProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var product models.Product
	initializers.DB.First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	initializers.DB.Model(&product).Updates(allProducts)
	c.JSON(http.StatusOK, gin.H{
		"product": product})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := initializers.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	initializers.DB.First(&product, id)

	if product.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found"})
		return
	}
	initializers.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{
		"message": "product deleted"})
}
