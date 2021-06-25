package main

import (
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/app"
	"github.com/RedisLabs-Field-Engineeringc/es-redis-proxy/handlers/healthcheck"
	"github.com/gin-gonic/gin"
	"github.com/shokunin/contrib/ginrus"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	router := gin.New()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true, "es-redis"))

	// Start routes
	router.GET("/health", healthcheck.HealthCheck)
	router.GET("/", app.AppRoot)

	// RUN rabit run
	router.Run() // listen and serve on 0.0.0.0:8080
}
