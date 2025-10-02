package middlewarex

import (
	"time"

	"SDT_ApiServices/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRouter returns a chi.Router with recommended middleware applied
func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	// ────── Core Middlewares ──────
	r.Use(middleware.RequestID)                 // Add unique request ID
	r.Use(middleware.RealIP)                    // Get client’s real IP
	r.Use(middleware.Logger)                    // Request logging
	r.Use(middleware.Recoverer)                 // Panic recovery
	r.Use(middleware.CleanPath)                 // Clean URL paths
	r.Use(middleware.Heartbeat("/ping"))        // Health check endpoint
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	// ────── Compression ──────
	r.Use(middleware.Compress(5)) // Compress responses (level 1–9)
	// 	1 → fastest compression, less size reduction
	// 9 → slowest, max compression

	// ────── CORS ──────
	r.Use(cors.Handler(config.MiddlewareRules)) // Centralized CORS config

	return r
}
