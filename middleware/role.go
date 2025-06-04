package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/models"
)

func CheckAdmin(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Forbidden"})
		return
	}
	if user.(models.User).Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Forbidden"})
	}
	c.Next()
}
