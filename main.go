package main

import (
	"log"

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
		&models.Address{},
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
		&models.Product{},
	)

	redisClient, err := initializers.RedisConnect()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	middleware.SetRedisClient(redisClient)
}

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Cors configurations
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Public auth routes
	r.POST("/auth/signup", controllers.CreateUser)
	r.POST("/auth/login", controllers.Login)

	// Protected routes
	auth := r.Group("/user")
	auth.Use(middleware.CheckAuth)
	auth.Use(middleware.RateLimiterMiddleware("default", middleware.GetRedisClient()))
	auth.GET("/profile", controllers.GetUserProfile)

	// Cart routes
	auth.GET("/cart", controllers.GetCart)
	auth.POST("/cart/items", controllers.AddToCart)
	auth.PUT("/cart/item/:id", controllers.DeleteCartItem)

	// Public product routes
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/products/:id", controllers.GetProductByID)

	// Protected product routes
	auth.POST("/products", controllers.CreateProduct)
	auth.PUT("/products/:id", controllers.UpdateProducts)
	auth.DELETE("/products/:id", controllers.DeleteProduct)

	// Order routes
	auth.POST("/orders", controllers.CreateOrder)
	auth.GET("/orders", controllers.GetUserOrder)
	auth.GET("/orders/:id", controllers.GetOrderByID)

	// Category routes
	r.GET("/categories", controllers.GetCategories)
	auth.POST("/categories", controllers.CreateCategory)

	// Payment routes
	auth.POST("/orders/:id/pay", controllers.CreateStripeCheckoutSession)
}
