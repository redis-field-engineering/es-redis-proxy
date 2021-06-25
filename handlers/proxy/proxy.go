package proxy

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {

	if c.Param("queryType") != "_search" {
		c.JSON(500, gin.H{
			"message": fmt.Sprintf("%s queries are not yet supported. Only _search", c.Param("queryType")),
		})
	} else {
		c.JSON(200, gin.H{
			"message": "This is proxypath",
		})
	}
}
