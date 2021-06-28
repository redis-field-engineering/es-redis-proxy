package main

import (
	"fmt"
	"time"

	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/gears"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/app"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/config"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/healthcheck"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/proxy"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/shokunin/contrib/ginrus"
	"github.com/sirupsen/logrus"
)

// Allows us to pass the Redis client handler to any handler functions
func APIMiddleware(r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redisConn", r)
		c.Next()
	}
}

func SetConfigCtx(config *config.ESProxyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("esProxyConfig", config)
		c.Next()
	}
}

func main() {

	config := config.GetConfig()

	router := gin.New()

	// Redis Setup
	redisClient := redis.NewClient(&redis.Options{
		Password:        config.RedisPassword,
		Addr:            fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 5000 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	router.Use(APIMiddleware(redisClient))
	router.Use(SetConfigCtx(config))

	gears.LoadGears(redisClient, "./gears/esSearch.py")

	logrus.SetFormatter(&logrus.JSONFormatter{})
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true, "es-redis"))

	// Grab our search
	// POST /instruments/_search HTTP/1.1
	router.POST("/:index/:queryType", proxy.Proxy)

	// Start routes
	router.GET("/health", healthcheck.HealthCheck)
	router.GET("/", app.AppRoot)

	// RUN rabit run
	router.Run() // listen and serve on 0.0.0.0:8080
}
