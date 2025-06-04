package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/loid-lab/e-commerce-api/utils"
)

func AddToCart(c *gin.Context) {
	var input struct {
		ProductID uint
		Quantity  int
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := utils.InvalidateKeys(initializers.RedisCLient, "carts")
	if err != nil {
		c.Error(err)
	}

	user, _ := c.Get("currentUser")

	cartItem := models.CartItem{
		UserID:    user.(models.User).ID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := initializers.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"car_item": cartItem})
}

func GetCart(c *gin.Context) {
	user, _ := c.Get("currentUser")
	cacheKey := fmt.Sprintf("cart:user:%d", user.(models.User).ID)

	var cartItem []models.CartItem
	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &cartItem)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"cart":   cartItem,
			"source": "cache"})
		return
	}

	initializers.DB.Where("user_id'=?", user.(models.User).ID).First(&cartItem)
	c.JSON(http.StatusOK, gin.H{
		"cart": cartItem})

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, cartItem, 5*time.Minute)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":   cartItem,
		"source": "db"})
}

func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("currentUser")

	var cartItem models.CartItem
	if err := initializers.DB.First(&cartItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart item not found"})
		return
	}

	if cartItem.UserID != user.(models.User).ID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access to cart item"})
		return
	}

	initializers.DB.Delete(&cartItem)

	err := utils.InvalidateKeys(initializers.RedisCLient, fmt.Sprintf("cart:user:%d", user.(models.User).ID))
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item deleted"})
}
