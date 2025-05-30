package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func CreateStripeCheckoutSession(c *gin.Context) {
	var input struct {
		OrderID uint `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("currentUser")

	var order models.Order
	if err := initializers.DB.First(&order, input.OrderID).Error; err != nil || order.UserID != user.(models.User).ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or not yours"})
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order #" + strconv.Itoa(int(order.ID))),
					},
					UnitAmount: stripe.Int64(int64(order.Total * 100)), // assuming Total is in float
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String("payment"),
		SuccessURL:        stripe.String("http://localhost:5173/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:         stripe.String("http://localhost:5173/cancel"),
		ClientReferenceID: stripe.String(strconv.Itoa(int(order.ID))),
	}

	s, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the session ID for reference
	payment := models.Payment{
		OrderID:  order.ID,
		Method:   "card",
		Provider: "Stripe",
		RefID:    s.ID,
		Status:   "pending",
	}
	initializers.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"checkout_url": s.URL})
}
