package models

// MongoAuthRequest represents MongoDB connection details
type MongoAuthRequest struct {
	URI      string `json:"uri" binding:"required" example:"mongodb://localhost:27017"`
	Database string `json:"database" binding:"required" example:"testdb"`
}

// MongoDynamicRequest represents dynamic MongoDB operations
type MongoDynamicRequest struct {
	DBConnection MongoAuthRequest       `json:"db_connection"`
	Action       string                 `json:"action" binding:"required,oneof=find insert update delete" example:"find"`
	Collection   string                 `json:"collection" binding:"required" example:"users"`
	Filter       map[string]interface{} `json:"filter" example:"{\"age\": {\"$gte\": 18}}"`
	Data         map[string]interface{} `json:"data" example:"{\"name\": \"Prashanth\", \"age\": 25}"`
	Page         int                    `json:"page" example:"1"`
	Limit        int                    `json:"limit" example:"10"`
	SortBy       string                 `json:"sort_by" example:"_id"`
	Order        string                 `json:"order" example:"asc"`
	IsMultiple   bool                   `json:"is_multiple"`
}
