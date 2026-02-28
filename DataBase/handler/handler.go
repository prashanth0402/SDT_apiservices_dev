package handler

import (
	"SDT_ApiServices/DataBase/models"
	"SDT_ApiServices/DataBase/service"
	"SDT_ApiServices/DataBase/validators"
	"errors"
	"strings"

	"gorm.io/gorm"
)

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
