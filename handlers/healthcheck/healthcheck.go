package healthcheck

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func HealthCheck(c *gin.Context) {
	var ctx = context.Background()
	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}
	fmt.Printf("%+v\n", redisConn)

	_, err := redisConn.Ping(ctx).Result()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Cannot ping Redis:",
			"error":   err,
		})
	} else {

		c.JSON(200, gin.H{
			"message": "OK",
		})
	}
}
