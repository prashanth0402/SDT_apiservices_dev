package gemini

import (
	models "SDT_ApiServices/Services/AI/Models"
	"SDT_ApiServices/utility"
	"context"
	"errors"
	"fmt"
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
	messages = append(messages, domain.NewTextMessage(domain.RoleUser, chatRequest.Prompt))
	result, err := p.GenerateMessage(ctx, messages, domain.WithMaxTokens(200))
	if utility.IsError(err) {
		return "", err
	}
	fmt.Println(result)
	if utility.IsEmpty(result.Content) {
		return "", errors.New("empty response from model")
	}
	return result.Content, nil
}
