package main

import (
	config "SDT_ApiServices/Config"
	services "SDT_ApiServices/Services"
	"SDT_ApiServices/middlewarex"
	"SDT_ApiServices/uihandlers"
	"fmt"
	"net/http"

	_ "SDT_ApiServices/docs"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           SDT API Services
// @version         1.0
// @description     Simple API service with a health check endpoint.
// @host            localhost:8080
// @BasePath        /
func main() {

	cfg := config.LoadConfig()

	// 🔹 Create Gin Router
	g := gin.Default()

	// 🔹 Load HTML templates
	g.LoadHTMLGlob("templates/**/*")

	// 🔹 Welcome Page (Internal Testing UI)

	g.GET("/", uihandlers.WelcomePage)

	// 🔹 Use your existing router inside Gin
	r := middlewarex.SetupRouter()

	// Existing APIs
	r.MethodFunc(http.MethodGet, "/getconnection", services.GetConnection)

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// 🔹 Mount old router into Gin
	g.Any("/api/*any", gin.WrapH(r))

	fmt.Println("Server running on http://localhost:" + cfg.Port)
	fmt.Println("Swagger: http://localhost:" + cfg.Port + "/api/swagger/index.html")

	// 🔹 Start Gin Server
	g.Run(":" + cfg.Port)
}
