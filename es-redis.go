package main

import (
	"fmt"
	"os"
	"time"

	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/app"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/healthcheck"
	"github.com/go-redis/redis/v8"

	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/proxy"
	"github.com/gin-gonic/gin"
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

func main() {

	var redisHost string
	var redisPort string
	var redisPassword string

	if len(os.Getenv("REDIS_HOST")) > 0 {
		redisHost = os.Getenv("REDIS_HOST")
	} else {
		redisHost = "localhost"
	}

	if len(os.Getenv("REDIS_PORT")) > 0 {
		redisPort = os.Getenv("REDIS_PORT")
	} else {
		redisPort = "6379"
	}

	router := gin.New()

	// Redis Setup
	redisClient := redis.NewClient(&redis.Options{
		Password:        redisPassword,
		Addr:            fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:              0,
		MinIdleConns:    1,                    // make sure there are at least this many connections
		MinRetryBackoff: 8 * time.Millisecond, //minimum amount of time to try and backupf
		MaxRetryBackoff: 5000 * time.Millisecond,
		MaxConnAge:      0,  //3 * time.Second this will cause everyone to reconnect every 3 seconds - 0 is keep open forever
		MaxRetries:      10, // retry 10 times : automatic reconnect if a proxy is killed
		IdleTimeout:     time.Second,
	})

	router.Use(APIMiddleware(redisClient))

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
