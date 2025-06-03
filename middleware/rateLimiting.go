package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loid-lab/e-commerce-api/models"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	limiter_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func GetRedisClient() *redis.Client {
	return redisClient
}

var globalRate = limiter.Rate{
	Period: 1 * time.Hour,
	Limit:  1000,
}

var routeNameMap = map[string]string{
	"/api/users": "users",
	"/api/items": "items",
}

var rateMapJSON = `{
	"default:users": "1000-H",
	"strict:users": "10-S",
	"default:items": "1000-H",
	"strict:items": "10-S"
}`

var key string

type RateConfig map[string]string

func retrieveRateConfig(mode, routeName string) (limiter.Rate, error) {
	var rateConfig RateConfig
	err := json.Unmarshal([]byte(rateMapJSON), &rateConfig)
	if err != nil {
		return limiter.Rate{}, err
	}

	key := fmt.Sprintf("%s:%s", mode, routeNameMap[routeName])
	rateStr, exists := rateConfig[key]
	if !exists {
		return limiter.Rate{}, fmt.Errorf("rate configuration not found for key %s", key)
	}

	rate, err := limiter.NewRateFromFormatted(rateStr)
	if err != nil {
		return limiter.Rate{}, err
	}

	return rate, nil
}

func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		ip, _, _ = net.SplitHostPort(c.Request.RemoteAddr)
	}
	return ip
}

func RateLimiterMiddleware(mode string, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		routePath := c.FullPath()
		rate, err := retrieveRateConfig(mode, routePath)
		if err != nil {
			rate = globalRate
		}

		if user, exists := c.Get("currentUser"); exists {
			key = fmt.Sprintf("user:%d", user.(models.User).ID)
		} else {
			key = fmt.Sprintf("ip:%s", getClientIP(c))
		}

		store, err := limiterRedis.NewStore(redisClient)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limiter store error"})
			return
		}

		lim := limiter.New(store, rate)
		limMiddleware := limiter_gin.NewMiddleware(lim)
		limMiddleware(c)
	}
}
