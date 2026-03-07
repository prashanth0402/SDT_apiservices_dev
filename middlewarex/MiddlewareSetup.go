package middlewarex

import (
	config "SDT_ApiServices/Config"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetupRouter returns a fully configured Gin router
func SetupRouter() *gin.Engine {

	r := gin.New()

	// Disable Console Color, you don't need console color when writing the logs to file.
	// gin.DisableConsoleColor()
	// Use the following code if you need to write the logs to file and console at the same time.
	f, _ := os.Create("Service.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// -----------------------------
	// Core Middlewares (Gin native)
	// -----------------------------
	r.Use(gin.Logger())   // logging
	r.Use(gin.Recovery()) // panic recovery
	// -----------------------------
	// Health Check (replacement for chi heartbeat)
	// -----------------------------
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// -----------------------------
	// Timeout middleware (custom)
	// -----------------------------
	r.Use(timeoutMiddleware(60 * time.Second))

	// -----------------------------
	// Compression
	// -----------------------------
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// -----------------------------
	// CORS
	// -----------------------------
	r.Use(cors.New(config.MiddlewareRules))

	return r
}

func timeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-time.After(timeout):
			c.JSON(504, gin.H{"error": "request timeout"})
			c.Abort()
		}
	}
}
