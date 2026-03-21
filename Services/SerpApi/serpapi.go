package serpapi

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/serpapi/serpapi-golang"
)

type GoogleSearchRequest struct {
	Query    string `json:"query" example:"Coffee"`
	Location string `json:"location" example:"Austin, Texas, United States"`
	Lang     string `json:"lang" example:"en"`
	Country  string `json:"country" example:"us"`
}

type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// GoogleSearch godoc
// @Summary Google Search API
// @Description Fetch Google search results using SerpAPI
// @Tags Search
// @Accept json
// @Produce json
// @Param request body GoogleSearchRequest true "Search Request"
// @Success 200 {array} SearchResult
// @Failure 400 {object} map[string]string
// @Router /googlesearch [post]
func GoogleSearch(c *gin.Context) {
	var req GoogleSearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	results, err := GoogleSearchResults(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func GoogleSearchResults(req GoogleSearchRequest) ([]SearchResult, error) {

	setting := serpapi.NewSerpApiClientSetting(os.Getenv("SERPAPI_KEY"))
	setting.Engine = "google"

	client := serpapi.NewClient(setting)

	params := map[string]string{
		"q":             req.Query,
		"location":      req.Location,
		"hl":            req.Lang,
		"gl":            req.Country,
		"google_domain": "google.com",
	}

	results, err := client.Search(params)
	if err != nil {
		return nil, err
	}

	organicResults := results["organic_results"].([]interface{})

	var response []SearchResult

	for _, item := range organicResults {
		data := item.(map[string]interface{})

		response = append(response, SearchResult{
			Title:   getString(data, "title"),
			Link:    getString(data, "link"),
			Snippet: getString(data, "snippet"),
		})
	}

	return response, nil
}

// helper
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
