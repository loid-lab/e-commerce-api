package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
)

func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, _ := c.Get("currentUser")
	order.UserID = user.(models.User).ID

	if err := initializers.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order})
}

func GetUserOrder(c *gin.Context) {
	var orders []models.Order
	if err := initializers.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"orders": orders})
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := initializers.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"order": order})
}
