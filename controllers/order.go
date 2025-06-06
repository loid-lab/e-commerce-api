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

func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := utils.InvalidateKeys(initializers.RedisCLient, "orders:all")
	if err != nil {
		c.Error(err)
	}

	user, _ := c.Get("currentUser")
	order.UserID = user.(models.User).ID

	if err := initializers.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not create order"})
		return
	}

	invoice := models.Invoice{
		InvoiceNumber: fmt.Sprintf("INV-%d", time.Now().Unix()),
		Date:          time.Now(),
		CustomerName:  user.(models.User).FullName,
		Items:         []models.InvoiceItem{},
	}

	for _, item := range order.Items {
		invoice.Items = append(invoice.Items, models.InvoiceItem{
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  item.TotalPrice,
		})
	}

	invoice.UserID = order.UserID

	if err := initializers.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create invoice"})
		return
	}

	email := models.EmailData{
		From:     "no-reply@example.com",
		To:       user.(models.User).Email,
		Subject:  "Your (company name) Invoice",
		HTMLBody: "<p>Thanks for your order! Your invoice is attached.</p>",
		SMTConfig: models.SMTConfig{
			SMTPHost: initializers.Env.SMTPHost,
			SMTPPort: initializers.Env.SMTPPort,
			SMTPUser: initializers.Env.SMTPUser,
			SMTPPass: initializers.Env.SMTPPass,
		},
	}

	go utils.GenerateSendInvoice(invoice, email)

	c.JSON(http.StatusOK, gin.H{
		"order": order})
}

func GetUserOrder(c *gin.Context) {
	var orders []models.Order
	cacheKey := "orders:all"

	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &orders)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"orders": orders,
			"source": "cache"})
		return
	}

	if err := initializers.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch products"})
		return
	}

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, orders, 10*time.Minute)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"source": "db",
	})
}

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	cacheKey := "order:" + id

	err := utils.GetJSON(initializers.RedisCLient, cacheKey, &order)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"order":  order,
			"source": "cache"})
		return
	}

	if err := initializers.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found"})
		return
	}

	err = utils.SetJSON(initializers.RedisCLient, cacheKey, order, 10*time.Minute)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"order":  order,
		"source": "db"})
}
