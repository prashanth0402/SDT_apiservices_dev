package models

type DBAuthRequest struct {
	Type     string `json:"type" binding:"required,oneof=mysql postgres" example:"mysql"`
	Host     string `json:"host" binding:"required" example:"localhost"`
	Port     string `json:"port" binding:"required,numeric" example:"3306"`
	Username string `json:"username" binding:"required" example:"root"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type DynamicRequest struct {
	Command   string                 `json:"command" example:"select"`
	TableName string                 `json:"table_name" example:"users"`
	Select    []string               `json:"select"` // columns
	Filters   []Filter               `json:"filters"`
	Data      map[string]interface{} `json:"data"`
	Page      int                    `json:"page" example:"1"`
	Limit     int                    `json:"limit" example:"10"`
	SortBy    string                 `json:"sort_by" example:"id"`
	Order     string                 `json:"order" example:"asc"`
}

type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // =, >, <, like, in
	Value    interface{} `json:"value"`
}
