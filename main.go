package main

import (
	config "SDT_ApiServices/Config"
	services "SDT_ApiServices/Services"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "SDT_ApiServices/docs" // MUST import the generated docs

	httpSwagger "github.com/swaggo/http-swagger" // Swagger UI handler
)

// @title           SDT API Services
// @version         1.0
// @description     Simple API service with a health check endpoint.
// @host            localhost:8080
// @BasePath        /
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // ✅ Global panic recovery
	// Load .env file

	cfg := config.LoadConfig()

	// ✅ Register your API endpoints
	r.MethodFunc(http.MethodGet, "/getconnection", services.GetConnection)

	// ✅ Swagger UI will be available at: http://localhost:8080/swagger/index.html
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:8080", cfg)
	fmt.Println("http://localhost:8080/swagger/index.html")
	http.ListenAndServe(":"+cfg.Port, r)
}
