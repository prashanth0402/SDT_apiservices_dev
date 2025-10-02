package config

import "github.com/go-chi/cors"

// Define your CORS rules as a variable, not const
var MiddlewareRules = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"}, // or []string{"*"}
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "userdevice"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300, // 5 minutes
}
