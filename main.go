package main

import (
	"SDT_ApiServices/DataBase"
	services "SDT_ApiServices/Services"
	"SDT_ApiServices/middlewarex"

	_ "SDT_ApiServices/docs"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

// func main() {

// 	cfg := config.LoadConfig()

// 	// 🔹 Create Gin Router
// 	g := gin.Default()

// 	// 🔹 Load HTML templates
// 	g.LoadHTMLGlob("templates/**/*")

// 	// 🔹 Welcome Page (Internal Testing UI)

// 	g.GET("/", uihandlers.WelcomePage)

// 	// 🔹 Use your existing router inside Gin
// 	r := middlewarex.SetupRouter()

// 	// Existing APIs
// 	r.MethodFunc(http.MethodGet, "/getconnection", services.GetConnection)

// 	// Swagger
// 	r.Get("/swagger/*", httpSwagger.WrapHandler)

// 	// 🔹 Mount old router into Gin
// 	g.Any("/api/*any", gin.WrapH(r))

// 	fmt.Println("Server running on http://localhost:" + cfg.Port)
// 	fmt.Println("Swagger: http://localhost:" + cfg.Port + "/api/swagger/index.html")

//		// 🔹 Start Gin Server
//		g.Run(":" + cfg.Port)
//	}
//
// @title           SDT API Services
// @version         1.0
// @description     Simple API service with a health check endpoint.
// @host            localhost:8080
// @BasePath        /
func main() {

	r := middlewarex.SetupRouter()

	r.GET("/getconnection", services.GetConnection)
	r.POST("/checkdb", DataBase.CheckDataBaseConnection)

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	r.Run(":8080")

}
