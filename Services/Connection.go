package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConnection godoc
// @Summary      Check API Connection
// @Description  Returns a simple JSON response to confirm that the API service is up.
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /getconnection [get]
func GetConnection(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Connection Established",
	})
}
