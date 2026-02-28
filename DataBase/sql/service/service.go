package service

import (
	"SDT_ApiServices/DataBase/SQL/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func HandleSelect(tx *gorm.DB, req *models.DynamicRequest) (interface{}, error) {

	var result []map[string]interface{}

	query := tx.Table(req.TableName)

	query = ApplyFilters(query, req.Filters)
	query = ApplySelect(query, req.Select)
	query = ApplySorting(query, req.SortBy, req.Order)
	query = ApplyPagination(query, req.Page, req.Limit)

	err := query.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ApplySelect(query *gorm.DB, columns []string) *gorm.DB {

	if len(columns) > 0 {
		query = query.Select(columns)
	}

	return query
}

func ApplyWhere(query *gorm.DB, where map[string]interface{}) *gorm.DB {

	for key, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	return query
}

func ApplyFilters(query *gorm.DB, filters []models.Filter) *gorm.DB {

	for _, f := range filters {

		switch strings.ToLower(f.Operator) {

		case "Eq":
			query = query.Where(fmt.Sprintf("%s = ?", f.Field), f.Value)

		case "gte":
			query = query.Where(fmt.Sprintf("%s > ?", f.Field), f.Value)

		case "lte":
			query = query.Where(fmt.Sprintf("%s < ?", f.Field), f.Value)

		case "like":
			query = query.Where(fmt.Sprintf("%s LIKE ?", f.Field), "%"+fmt.Sprint(f.Value)+"%")
		case "notlike":
			query = query.Where(fmt.Sprintf("%s NOT LIKE ?", f.Field), "%"+fmt.Sprint(f.Value)+"%")

		case "in":
			query = query.Where(fmt.Sprintf("%s IN ?", f.Field), f.Value)
		case "notin":
			query = query.Where(fmt.Sprintf("%s NOT IN ?", f.Field), f.Value)
		}
	}
	return query
}

func ApplyPagination(query *gorm.DB, page int, limit int) *gorm.DB {

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	return query.Offset(offset).Limit(limit)
}

func HandleCreate(tx *gorm.DB, req *models.DynamicRequest) (interface{}, error) {

	if req.Data == nil {
		return nil, errors.New("data required")
	}

	err := tx.Table(req.TableName).Create(req.Data).Error
	if err != nil {
		return nil, err
	}

	return "record created successfully", nil
}

func ApplySorting(query *gorm.DB, sortBy string, order string) *gorm.DB {

	if sortBy == "" {
		return query
	}

	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return query.Order(sortBy + " " + order)
}

func HandleDelete(tx *gorm.DB, req *models.DynamicRequest) (interface{}, error) {

	query := tx.Table(req.TableName)
	query = ApplyFilters(query, req.Filters)

	err := query.Delete(nil).Error
	if err != nil {
		return nil, err
	}

	return "record deleted successfully", nil
}

func HandleUpdate(tx *gorm.DB, req *models.DynamicRequest) (interface{}, error) {

	if req.Data == nil {
		return nil, errors.New("data required")
	}

	query := tx.Table(req.TableName)
	query = ApplyFilters(query, req.Filters)

	err := query.Updates(req.Data).Error
	if err != nil {
		return nil, err
	}

	return "record updated successfully", nil
}
