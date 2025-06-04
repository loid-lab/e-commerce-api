package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/loid-lab/e-commerce-api/utils"
)

func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := c.Get("currentUser")
	category.CreatedBy = user.(models.User).ID

	if err := initializers.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"category": category})

	err := utils.InvalidateKeys(initializers.RedisCLient, "categories:all")
	if err != nil {
		c.Error(err)
	}
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	cacheKey := "categories:all"

	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &categories)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"categories": categories,
			"sources":    "cache"})
	}

	if err := initializers.DB.Find(&categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch Categories"})
		return
	}

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, categories, 5*time.Minute)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"source":     "db"})
}
