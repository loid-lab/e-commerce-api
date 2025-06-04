package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/controllers"
	"github.com/loid-lab/e-commerce-api/initializers"
	"github.com/loid-lab/e-commerce-api/middleware"
	"github.com/loid-lab/e-commerce-api/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.Payment{},
		&models.Invoice{},
		&models.InvoiceItem{},
	)

	redisClient, err := initializers.RedisConnect()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	middleware.SetRedisClient(redisClient)

	initializers.ConnectCloudinary()
}

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Cors configurations
	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/auth/signup", controllers.CreateUser)
	r.POST("/auth/login", controllers.Login)
	r.POST("/webhooks/stripe", controllers.StripeWebhook)

	// Admin routes
	admin := r.Group("/admin")
	admin.Use(middleware.CheckAuth, middleware.CheckAdmin)
	{
		admin.GET("/metrics/sales", controllers.GetSalesMetrics)
		admin.GET("/orders/stats", controllers.GetOrderStats)
		admin.GET("/invoices", controllers.GetAllInvoices)
	}

	// Authenticated user routes
	auth := r.Group("/user")
	auth.Use(middleware.CheckAuth)
	auth.Use(middleware.RateLimiterMiddleware("default", middleware.GetRedisClient()))
	{
		auth.GET("/profile", controllers.GetUserProfile)

		auth.GET("/cart", controllers.GetCart)
		auth.POST("/cart/items", controllers.AddToCart)
		auth.PUT("/cart/item/:id", controllers.DeleteCartItem)

		auth.POST("/orders", controllers.CreateOrder)
		auth.GET("/orders", controllers.GetUserOrder)
		auth.GET("/orders/:id", controllers.GetOrderByID)

		auth.POST("/orders/:id/pay", controllers.CreateStripeCheckoutSession)

		auth.POST("/products", middleware.CheckAdmin, controllers.CreateProduct)
		auth.PUT("/products/:id", middleware.CheckAdmin, controllers.UpdateProducts)
		auth.DELETE("/products/:id", middleware.CheckAdmin, controllers.DeleteProduct)

		auth.POST("/categories", middleware.CheckAdmin, controllers.CreateCategory)
	}

	// Public product and category routes
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/products/:id", controllers.GetProductByID)
	r.GET("/categories", controllers.GetCategories)

	r.Run()
}
