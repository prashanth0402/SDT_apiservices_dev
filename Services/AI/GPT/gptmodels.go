package gpt

import (
	"SDT_ApiServices/utility"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lexlapax/go-llms/pkg/llm/domain"
	"github.com/lexlapax/go-llms/pkg/llm/provider"

	"github.com/gin-gonic/gin"
)

type ChatBotRequest struct {
	OpenAI_API_Key string `json:"openAi_ApI_key" example:"sk-1234567890abcdef"`
	Prompt         string `json:"prompt" example:"Explain Golang in simple terms"`
	Agent          string `json:"agent" example:"You are a helpful assistant"`
	Model          string `json:"model" example:"gpt-4o-mini"`
}

// GPTHandler godoc
// @Summary Chat with GPT
// @Description Send prompt and get AI response
// @Tags GPT
// @Accept json
// @Produce json
// @Param request body ChatBotRequest true "Chat Request"
// @Failure 400 {object} map[string]string
// @Router /gpthandler [post]
func GPTHandler(c *gin.Context) {
	var req ChatBotRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Timeout context (important)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	response, err := ChatwithGPT(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

// ChatwithGPT handles LLM interaction
func ChatwithGPT(ctx context.Context, chatRequest ChatBotRequest) (string, error) {

	if utility.IsEmpty(chatRequest.OpenAI_API_Key) {
		return "", errors.New("API key is required")
	}

	// Default model
	model := chatRequest.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	// Create provider
	p := provider.NewOpenAIProvider(chatRequest.OpenAI_API_Key, model)

	// System prompt
	// systemPrompt := "You are a helpful assistant."
	// if chatRequest.Agent != "" {
	// 	systemPrompt = chatRequest.Agent
	// }

	// Get user input
	if chatRequest.Prompt == "" {
		return "", errors.New("input prompt is required")
	}

	// Prepare messages (THIS replaces state)
	// messages := []domain.Message{
	// 	domain.NewTextMessage(domain.RoleSystem,systemPrompt),domain.NewTextMessage(domain.RoleUser,userInput)
	// }

	var messages []domain.Message
	messages = append(messages, domain.NewTextMessage(domain.RoleUser, chatRequest.Prompt))
	// Call LLM (CORRECT WAY ✅)
	result, err := p.GenerateMessage(ctx, messages,
		domain.WithTemperature(0.7),
		domain.WithMaxTokens(200),
	)
	if err != nil {
		return "", err
	}

	// Extract response safely
	if result.Content == "" {
		return "", errors.New("empty response from model")
	}

	return result.Content, nil
}
