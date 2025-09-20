package services

import (
	"SDT_ApiServices/common"
	"fmt"
	"net/http"
)

// GetConnection godoc
// @Summary      Check API Connection
// @Description  Returns a simple JSON response to confirm that the API service is up.
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "Returns connection status"
// @Failure      405  {object}  map[string]string  "Method Not Allowed"
// @Router       /getconnection [get]
func GetConnection(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", common.AllowOrgin)
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message": "Connection Established"}`)
}
