package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/loid-lab/e-commerce-api/utils"
)

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := utils.InvalidateKeys(initializers.RedisCLient, "products:all")
	if err != nil {
		c.Error(err)
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
	cacheKey := "products:all"

	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &products)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"products": products,
			"source":   "cache"})
		return
	}
	if err := initializers.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, products, 10*time.Minute)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK,
		gin.H{
			"products": products,
			"source":   "db"})
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

	err := utils.InvalidateKeys(initializers.RedisCLient, "product:all", "product:"+id)
	if err != nil {
		c.Error(err)
	}
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	cacheKey := "product:" + id

	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &product)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"product": product,
			"source":  "cache"})
	}

	if err := initializers.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{
			"product": product,
			"source":  "db"})

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, product, 10*time.Minute)
	if err != nil {
		c.Error(err)
	}
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

	err := utils.InvalidateKeys(initializers.RedisCLient, "products:all", "product:"+id)
	if err != nil {
		c.Error(err)
	}
}
