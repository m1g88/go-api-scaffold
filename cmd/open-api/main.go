package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1g88/go-api-scaffold/internal/healthcheck"
	"github.com/m1g88/go-api-scaffold/pkg/config"
	"github.com/m1g88/go-api-scaffold/pkg/httpserver"
	"github.com/m1g88/go-api-scaffold/pkg/logger"
)

func main() {
	config.Init()
	r := gin.New()

	log := logger.New()
	r.Use(log.MiddlewareLogger())
	r.Use(log.RecoveryWithZap())
	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example when panic happen.
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})

	healthcheck.SetHealthCheckRoute(r)

	server := httpserver.NewGinServer(
		httpserver.WithTimeout(30 * time.Second))
	server.Handler = r

	server.Run()
}
