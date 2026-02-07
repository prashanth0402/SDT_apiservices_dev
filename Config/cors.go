package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

// Define your CORS rules as a variable, not const
var MiddlewareRules = cors.Config{
	AllowOrigins:     []string{"http://*", "https://*"}, // or []string{"*"}
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "userdevice"},
	ExposeHeaders:    []string{"Link"},
	AllowCredentials: true,
	MaxAge:           5 * time.Minute,
}
