package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/models"
)

func GetUserProfile(c *gin.Context) {
	user, exist := c.Get("currentUser")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		c.Abort()
		return
	}
	currentUser := user.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    currentUser.ID,
			"email": currentUser.Email,
		},
	})
}
