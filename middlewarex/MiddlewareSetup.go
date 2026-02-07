// package middlewarex

// import (
// 	config "SDT_ApiServices/Config"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/go-chi/cors"
// )

// // SetupRouter returns a chi.Router with recommended middleware applied
// func SetupRouter() *chi.Mux {
// 	r := gin.Default()

// 	// ────── Core Middlewares ──────
// 	r.Use(middleware.RequestID())                 // Add unique request ID
// 	r.Use(middleware.RealIP)                    // Get client’s real IP
// 	r.Use(middleware.Logger)                    // Request logging
// 	r.Use(middleware.Recoverer)                 // Panic recovery
// 	r.Use(middleware.CleanPath)                 // Clean URL paths
// 	r.Use(middleware.Heartbeat("/ping"))        // Health check endpoint
// 	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

// 	// ────── Compression ──────
// 	r.Use(middleware.Compress(5)) // Compress responses (level 1–9)
// 	// 	1 → fastest compression, less size reduction
// 	// 9 → slowest, max compression

// 	// ────── CORS ──────
// 	r.Use(cors.Handler(config.MiddlewareRules)) // Centralized CORS config

// 	return r
// }

package middlewarex

import (
	config "SDT_ApiServices/Config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetupRouter returns a fully configured Gin router
func SetupRouter() *gin.Engine {

	r := gin.New()

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
