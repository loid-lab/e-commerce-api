package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
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

	var cart []models.Cart

	initializers.DB.Where("user_id'=?", user.(models.User).ID).First(&cart)
	c.JSON(http.StatusOK, gin.H{
		"cart": cart})
}

func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")

	var carts models.Cart
	initializers.DB.First(&carts, id)

	if carts.UserID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cart not found"})
		return
	}
	initializers.DB.Delete(&carts)
	c.JSON(http.StatusOK, gin.H{
		"message": "Cart deleted"})
}
