package proxy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TriggerResponse struct {
	Result      interface{} `json:"result"`
	TTL         int         `json:"ttl"`
	ExitCode    int         `json:"exit_code"`
	CacheStatus string      `json:"cache_status"`
}

func Proxy(c *gin.Context) {

	var tr TriggerResponse

	var ctx = context.Background()

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	jsonQueryHash := sha256.New()
	jsonQueryHash.Write(jsonData)

	esProxyConfig, cok := c.MustGet("esProxyConfig").(*config.ESProxyConfig)
	if !cok {
		c.JSON(500, gin.H{
			"message": "Cannot get configuration",
		})
	}

	redisConn, ok := c.MustGet("redisConn").(*redis.Client)
	if !ok {
		c.JSON(500, gin.H{
			"message": "Cannot get redisConn",
		})
	}

	if c.Param("queryType") != "_search" {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("%s queries are not yet supported. Only _search", c.Param("queryType")),
		})
	}
	trigger, terr := redisConn.Do(
		ctx,
		"RG.TRIGGER",
		"es-search",
		fmt.Sprintf("%x", (jsonQueryHash.Sum(nil))),
		c.Param("index"),
		jsonData,
		esProxyConfig.RedisTTL,
		esProxyConfig.ReCacheInterval,
	).Result()
	if terr != nil {
		c.JSON(500, terr)
	} else {
		json.Unmarshal([]byte(trigger.([]interface{})[0].(string)), &tr)
		c.Header("X-Cache-TTL", fmt.Sprintf("%d", tr.TTL))
		c.Header("X-Cache-Status", tr.CacheStatus)
		c.Header("X-Cache-Upstream", fmt.Sprintf("%d", tr.ExitCode))
		httpStatus := 503
		if tr.ExitCode < 1 {
			httpStatus = 200
		}
		c.JSON(httpStatus, tr.Result)
	}
}
