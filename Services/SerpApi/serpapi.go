package serpapi

import (
	"net/http"
	"os"
	"strconv"

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

// JobSearchRequest is the query used against SerpApi's Google Jobs engine,
// which returns real, live job listings aggregated from LinkedIn, Indeed,
// Naukri, company career pages, etc. — no site-specific scraping required.
type JobSearchRequest struct {
	Query    string `json:"query" example:"Software Engineer"`
	Location string `json:"location" example:"Bangalore, India"`
	Lang     string `json:"lang" example:"en"`
	Country  string `json:"country" example:"in"`
	Page     int    `json:"page" example:"0"` // 0-based; each page is 10 results
}

// JobResult is a normalized subset of one entry from SerpApi's
// "jobs_results" array.
type JobResult struct {
	JobID        string   `json:"job_id"`
	Title        string   `json:"title"`
	Company      string   `json:"company_name"`
	Location     string   `json:"location"`
	Via          string   `json:"via"`
	Description  string   `json:"description"`
	PostedAt     string   `json:"posted_at,omitempty"`
	ScheduleType string   `json:"schedule_type,omitempty"`
	Salary       string   `json:"salary,omitempty"`
	ApplyLink    string   `json:"apply_link,omitempty"`
	Extensions   []string `json:"extensions,omitempty"`
}

// GoogleJobsSearch godoc
// @Summary Google Jobs Search API
// @Description Fetch real, live job listings using SerpAPI's Google Jobs engine
// @Tags Search
// @Accept json
// @Produce json
// @Param request body JobSearchRequest true "Job Search Request"
// @Success 200 {array} JobResult
// @Failure 400 {object} map[string]string
// @Router /googlejobssearch [post]
func GoogleJobsSearch(c *gin.Context) {
	var req JobSearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	results, err := GoogleJobsSearchResults(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GoogleJobsSearchResults queries SerpApi's google_jobs engine and normalizes
// each entry in "jobs_results" into a JobResult.
func GoogleJobsSearchResults(req JobSearchRequest) ([]JobResult, error) {

	setting := serpapi.NewSerpApiClientSetting(os.Getenv("SERPAPI_KEY"))
	setting.Engine = "google_jobs"

	client := serpapi.NewClient(setting)

	params := map[string]string{
		"q":             req.Query,
		"location":      req.Location,
		"hl":            req.Lang,
		"gl":            req.Country,
		"google_domain": "google.com",
	}
	if req.Page > 0 {
		params["start"] = strconv.Itoa(req.Page * 10)
	}

	results, err := client.Search(params)
	if err != nil {
		return nil, err
	}

	jobsRaw, ok := results["jobs_results"].([]interface{})
	if !ok {
		return []JobResult{}, nil
	}

	jobs := make([]JobResult, 0, len(jobsRaw))
	for _, item := range jobsRaw {
		data, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		job := JobResult{
			JobID:       getString(data, "job_id"),
			Title:       getString(data, "title"),
			Company:     getString(data, "company_name"),
			Location:    getString(data, "location"),
			Via:         getString(data, "via"),
			Description: getString(data, "description"),
		}

		if detected, ok := data["detected_extensions"].(map[string]interface{}); ok {
			job.PostedAt = getString(detected, "posted_at")
			job.ScheduleType = getString(detected, "schedule_type")
			job.Salary = getString(detected, "salary")
		}

		if exts, ok := data["extensions"].([]interface{}); ok {
			for _, e := range exts {
				if s, ok := e.(string); ok {
					job.Extensions = append(job.Extensions, s)
				}
			}
		}

		if applyOptions, ok := data["apply_options"].([]interface{}); ok && len(applyOptions) > 0 {
			if first, ok := applyOptions[0].(map[string]interface{}); ok {
				job.ApplyLink = getString(first, "link")
			}
		} else if relatedLinks, ok := data["related_links"].([]interface{}); ok && len(relatedLinks) > 0 {
			if first, ok := relatedLinks[0].(map[string]interface{}); ok {
				job.ApplyLink = getString(first, "link")
			}
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
