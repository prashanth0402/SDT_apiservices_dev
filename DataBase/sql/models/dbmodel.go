package models

// DBAuthRequest represents database authentication details
type DBAuthRequest struct {
	Type     string `json:"type" binding:"required,oneof=mysql postgres" example:"mysql"`
	Host     string `json:"host" binding:"required" example:"localhost"`
	Port     string `json:"port" binding:"required,numeric" example:"3306"`
	Username string `json:"username" binding:"required" example:"root"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// Filter represents dynamic query filter condition
type Filter struct {
	Field    string      `json:"field" example:"age"`
	Operator string      `json:"operator" example:"gte"` // eq, gte, lte, like, in, notlike, notin
	Value    interface{} `json:"value" example:"25"`
}

// DynamicRequest represents the dynamic DB operation request
type DynamicRequest struct {
	DBConnection DBAuthRequest          `json:"db_connection"`
	Command      string                 `json:"command" binding:"required,oneof=create select update delete" example:"select"`
	TableName    string                 `json:"table_name" binding:"required" example:"users"`
	Select       []string               `json:"select" example:"id,name,email"`
	Filters      []Filter               `json:"filters"`
	Data         map[string]interface{} `json:"data"`
	Page         int                    `json:"page" example:"1"`
	Limit        int                    `json:"limit" example:"10"`
	SortBy       string                 `json:"sort_by" example:"id"`
	Order        string                 `json:"order" example:"asc"`
}
