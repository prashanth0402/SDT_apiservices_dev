package validators

import (
	"github.com/prashanth0402/SDT_apiservices_dev/DataBase/sql/models"
	"errors"
)

var allowedTables = map[string]bool{
	"users":  true,
	"orders": true,
}

func ValidateRequest(req *models.DynamicRequest) error {

	if req.Command == "" {
		return errors.New("command required")
	}

	if req.TableName == "" {
		return errors.New("table name required")
	}

	// if !allowedTables[req.TableName] {
	// 	return errors.New("table not allowed")
	// }

	return nil
}
