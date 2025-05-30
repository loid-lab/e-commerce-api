package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
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
}

func GetCategories(c *gin.Context) {
	var categories []models.Category

	if err := initializers.DB.Find(&categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch Categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories})
}
