package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

func StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	payload, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusServiceUnavailable, "Error reading request body")
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		c.String(http.StatusBadRequest,
			"Signature verification failed")
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.String(http.StatusBadRequest, "Failed to parse event data")
			return
		}

		orderID := session.ClientReferenceID

		var payment models.Payment
		if err := initializers.DB.Where("ref_id=?", session.ID).First(&payment).Error; err == nil {
			payment.Status = "paid"
			initializers.DB.Save(&payment)
		}

		if orderID != "" {
			var order models.Order
			if err := initializers.DB.First(&order, orderID).Error; err == nil {
				order.Status = "paid"
				initializers.DB.Save(&order)
			}
		}

		c.Status(http.StatusOK)
	default:
		c.Status(http.StatusOK)
	}
}
