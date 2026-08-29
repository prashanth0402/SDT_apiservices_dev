package gemini

import (
	"context"
	"errors"
	models "github.com/prashanth0402/SDT_apiservices_dev/Services/AI/Models"
	"github.com/prashanth0402/SDT_apiservices_dev/utility"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lexlapax/go-llms/pkg/llm/domain"
	"github.com/lexlapax/go-llms/pkg/llm/provider"
)

// ChatWithGemini godoc
// @Summary Chat with Gemini AI
// @Description Send prompt to Gemini and get response
// @Tags Gemini
// @Accept json
// @Produce json
// @Param request body models.ChatBotRequest true "ChatBotRequest"
// @Failure 400 {object} map[string]string
// @Router /geminihandler [post]
func GemniHandler(c *gin.Context) {
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

	response, err := ChatwithGemini(ctx, req)
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

func ChatwithGemini(ctx context.Context, chatRequest models.ChatBotRequest) (string, error) {
	if utility.IsEmpty(chatRequest.GeminiAI_API_Key) {
		return "", errors.New("API key is required")
	}
	p := provider.NewGeminiProvider(chatRequest.GeminiAI_API_Key, chatRequest.Model)
	var messages []domain.Message
	if chatRequest.Agent != "" {
		messages = append(messages, domain.NewTextMessage(domain.RoleSystem, chatRequest.Agent))
	}
	messages = append(messages, domain.NewTextMessage(domain.RoleUser, chatRequest.Prompt))

	maxTokens := 200
	if chatRequest.MaxTokens > 0 {
		maxTokens = chatRequest.MaxTokens
	}
	temperature := 0.7
	if chatRequest.Temperature > 0 {
		temperature = chatRequest.Temperature
	}
	result, err := p.GenerateMessage(ctx, messages,
		domain.WithMaxTokens(maxTokens),
		domain.WithTemperature(temperature),
	)
	if utility.IsError(err) {
		return "", err
	}
	if utility.IsEmpty(result.Content) {
		return "", errors.New("empty response from model")
	}
	return result.Content, nil
}
