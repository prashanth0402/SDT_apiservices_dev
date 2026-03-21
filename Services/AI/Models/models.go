package models

type ChatBotRequest struct {
	OpenAI_API_Key   string `json:"openAi_ApI_key" example:"sk-1234567890abcdef"`
	GeminiAI_API_Key string `json:"geminiAi_ApI_key" example:"sk-1234567890abcdef"`
	Prompt           string `json:"prompt" example:"Explain Golang in simple terms"`
	Agent            string `json:"agent" example:"You are a helpful assistant"`
	Model            string `json:"model" example:"gpt-4o-mini"`
}
