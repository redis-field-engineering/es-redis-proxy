package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"crypto/sha256"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
)

var (
	esResponse map[string]interface{}
)

func Proxy(c *gin.Context) {

	var ctx = context.Background()
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
	} else {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		}
		jsonQueryHash := sha256.New()
		jsonQueryHash.Write(jsonData)
		res, rerr := redisConn.Get(ctx,
			fmt.Sprintf("%x", (jsonQueryHash.Sum(nil))),
		).Result()
		if rerr != nil {
			// Go ahed and fetch this from elasticsearch
			es, eserr := elasticsearch.NewDefaultClient()
			if eserr != nil {
				fmt.Println("ES Error: ", eserr)
			}
			sres, serr := es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(c.Param("index")),
				es.Search.WithBody(bytes.NewReader(jsonData)),
				es.Search.WithTrackTotalHits(true),
				es.Search.WithPretty(),
			)
			if serr != nil {
				fmt.Println("Search error: ", serr)
			}
			defer sres.Body.Close()
			if sres.IsError() {
				var e map[string]interface{}
				if err := json.NewDecoder(sres.Body).Decode(&e); err != nil {
					fmt.Sprintf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and error information.
					fmt.Sprintf("[%s] %s: %s",
						sres.Status(),
						e["error"].(map[string]interface{})["type"],
						e["error"].(map[string]interface{})["reason"],
					)
				}
			}
			if err := json.NewDecoder(sres.Body).Decode(&esResponse); err != nil {
				fmt.Printf("Error parsing the response body: %s", err)
			}
			jsonString, _ := json.Marshal(esResponse)
			_, xerr := redisConn.SetNX(
				ctx,
				fmt.Sprintf("%x", (jsonQueryHash.Sum(nil))),
				jsonString,
				10*time.Second,
			).Result()
			if xerr != nil {
				fmt.Printf("ERROR SETTING %+v\n", xerr)
			}
			c.Header("X-Cache", "0")
			c.JSON(200, esResponse)
		} else {
			if err := json.NewDecoder(strings.NewReader(res)).Decode(&esResponse); err != nil {
				fmt.Printf("Error parsing the response body: %s", err)
			}
			c.Header("X-Cache", "1")
			c.JSON(200, esResponse)
		}
	}
}
