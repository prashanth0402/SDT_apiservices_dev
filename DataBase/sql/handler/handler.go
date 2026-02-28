package handler

import (
	DataBase "SDT_ApiServices/DataBase/SQL"
	"SDT_ApiServices/DataBase/SQL/models"
	"SDT_ApiServices/DataBase/SQL/service"
	"SDT_ApiServices/DataBase/SQL/validators"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExecuteSQlQuery godoc
// @Summary Execute Multi-Tenant Dynamic Query
// @Description DB-as-a-Service dynamic CRUD endpoint supporting filters, operators, pagination and sorting
// @Tags Tenant Query Engine
// @Accept json
// @Produce json
// @Param request body models.DynamicRequest true "Dynamic Query Request"
// @Success 200 {object} interface{} "Successful response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /excecutesqlquery [post]

func ExecuteSQLQuery(c *gin.Context) {

	var req models.DynamicRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dsn, err := DataBase.BuildDSN(req.DBConnection)

	db, err := DataBase.ConnectDB(req.DBConnection.Type, dsn)

	result, err := ExecuteDynamic(&db, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)

}

func ExecuteDynamic(db *gorm.DB, req *models.DynamicRequest) (interface{}, error) {

	if err := validators.ValidateRequest(req); err != nil {
		return nil, err
	}

	var result interface{}
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {

		switch strings.ToLower(req.Command) {

		case "Insert":
			result, err = service.HandleCreate(tx, req)

		case "select":
			result, err = service.HandleSelect(tx, req)

		case "update":
			result, err = service.HandleUpdate(tx, req)

		case "delete":
			result, err = service.HandleDelete(tx, req)

		default:
			return errors.New("invalid command")
		}

		return err
	})

	return result, err
}
