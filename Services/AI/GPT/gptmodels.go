package gpt

import (
	models "SDT_ApiServices/Services/AI/Models"
	"SDT_ApiServices/utility"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lexlapax/go-llms/pkg/llm/domain"
	"github.com/lexlapax/go-llms/pkg/llm/provider"

	"github.com/gin-gonic/gin"
)

// GPTHandler godoc
// @Summary Chat with GPT
// @Description Send prompt and get AI response
// @Tags GPT
// @Accept json
// @Produce json
// @Param request body models.ChatBotRequest true "ChatBotRequest"
// @Failure 400 {object} map[string]string
// @Router /gpthandler [post]
func GPTHandler(c *gin.Context) {
	var req models.ChatBotRequest

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
func ChatwithGPT(ctx context.Context, chatRequest models.ChatBotRequest) (string, error) {

	if utility.IsEmpty(chatRequest.OpenAI_API_Key) {
		return "", errors.New("API key is required")
	}

	// Default model
	if utility.IsEmpty(chatRequest.Model) {
		chatRequest.Model = "gpt-4o-mini"
	}

	// Create provider
	p := provider.NewOpenAIProvider(chatRequest.OpenAI_API_Key, chatRequest.Model)

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
	if utility.IsError(err) {
		return "", err
	}

	// Extract response safely
	if result.Content == "" {
		return "", errors.New("empty response from model")
	}

	return result.Content, nil
}
