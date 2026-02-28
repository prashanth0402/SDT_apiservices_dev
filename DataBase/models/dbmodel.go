package models

type DBAuthRequest struct {
	Type     string `json:"type" binding:"required,oneof=mysql postgres" example:"mysql"`
	Host     string `json:"host" binding:"required" example:"localhost"`
	Port     string `json:"port" binding:"required,numeric" example:"3306"`
	Username string `json:"username" binding:"required" example:"root"`
	Password string `json:"password" binding:"required" example:"123456"`
}
