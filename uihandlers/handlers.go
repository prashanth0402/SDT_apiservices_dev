package uihandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title":   "SDT API Services",
		"message": "Welcome to SDT Internal Testing UI 🚀",
	})
}
