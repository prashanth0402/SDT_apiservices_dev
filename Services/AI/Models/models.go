package models

type ChatBotRequest struct {
	OpenAI_API_Key   string `json:"openAi_ApI_key" example:"sk-1234567890abcdef"`
	GeminiAI_API_Key string `json:"geminiAi_ApI_key" example:"sk-1234567890abcdef"`
	Prompt           string `json:"prompt" example:"Explain Golang in simple terms"`
	Agent            string `json:"agent" example:"You are a helpful assistant"`
	Model            string `json:"model" example:"gpt-4o-mini"`
	// Optional generation controls. Zero means "use the service default"
	// (MaxTokens 200, Temperature 0.7) so existing callers are unaffected.
	MaxTokens   int     `json:"maxTokens,omitempty" example:"800"`
	Temperature float64 `json:"temperature,omitempty" example:"0.2"`
}
