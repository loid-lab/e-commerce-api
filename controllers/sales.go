package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
)

func GetSalesMetrics(c *gin.Context) {
	var totalSales float64
	var orderCount int64

	initializers.DB.Model(&models.Order{}).Where("status = ?", "paid").Count(&orderCount)
	initializers.DB.Model(&models.Order{}).Select("SUM(total)").Where("status = ?", "paid").Scan(&totalSales)
}

func GetOrderStats(c *gin.Context) {
	var total int64
	var paidCount int64

	if err := initializers.DB.Model(&models.Order{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count paid orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_orders": total,
		"paid_orders":  paidCount})
}
