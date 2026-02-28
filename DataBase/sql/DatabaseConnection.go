package DataBase

import (
	"SDT_ApiServices/DataBase/SQL/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CheckDataBaseConnection godoc
// @Summary      Check Database Server Connection
// @Description  Checks only server connectivity (host + port + credentials). Does NOT verify database existence.
// @Tags         Database
// @Accept       json
// @Produce      json
// @Param        request  body      models.DBAuthRequest  true  "Database Authentication Details"
// @Success      200  {object}  map[string]interface{}  "Connection successful"
// @Failure      400  {object}  map[string]interface{}  "Invalid request body"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /checkdb [post]
func CheckDataBaseConnection(c *gin.Context) {

	var req models.DBAuthRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// ✅ Build DSN
	dsn, err := BuildDSN(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// ✅ Connect
	db, err := ConnectDB(req.Type, dsn)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Connection failed",
			"error":   err.Error(),
		})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get SQL instance",
		})
		return
	}

	// ✅ Ping only server
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Server not reachable",
			"error":   err.Error(),
		})
		return
	}

	defer sqlDB.Close()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database server is reachable",
	})
}

func BuildDSN(req models.DBAuthRequest) (string, error) {

	switch req.Type {

	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
			req.Username,
			req.Password,
			req.Host,
			req.Port,
		), nil

	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s port=%s dbname=postgres sslmode=disable",
			req.Host,
			req.Username,
			req.Password,
			req.Port,
		), nil

	default:
		return "", fmt.Errorf("unsupported database type")
	}
}

func ConnectDB(dbType, dsn string) (*gorm.DB, error) {

	switch dbType {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// case "sqlite":
	// 	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type")
	}
}
