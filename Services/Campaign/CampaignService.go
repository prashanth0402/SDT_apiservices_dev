package campaign

import (
	email "SDT_ApiServices/Services/Campaign/Email"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
)

// CampaignRequest defines request payload
type CampaignRequest struct {
	// Channel type: email, whatsapp, telegram
	Channel string `json:"channel" validate:"required,oneof=email whatsapp telegram"`
	// Emails to send (for email campaigns)
	Emails []email.EmailInput `json:"emails" validate:"required,dive"`
}

// Messenger interface for campaign channels
type Messenger interface {
	SendEmail(e email.EmailInput) error
}

var validate = validator.New()

// CampaignHandler handles campaign requests
// @Summary Start a messaging campaign
// @Description Send campaign messages via email/WhatsApp/Telegram
// @Tags Campaign
// @Accept json
// @Produce json
// @Param campaign body CampaignRequest true "Campaign Request"
// @Success 202 {string} string "Campaign started"
// @Failure 400 {string} string "Invalid request"
// @Router /campaign [post]
func CampaignHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// JSON validation
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	var messenger Messenger
	switch req.Channel {
	case "email":
		messenger = &email.EmailSender{}
		go SendEmailCampaign(messenger, req, len(req.Emails)) // Adapter to match interface
	// case "whatsapp":
	// 	messenger = channels.WhatsAppSender{}
	// case "telegram":
	// 	messenger = channels.TelegramSender{}
	default:
		http.Error(w, "Unsupported channel", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Campaign started"))
}

// SendEmailCampaign sends emails concurrently using worker pool
func SendEmailCampaign(messenger Messenger, req CampaignRequest, concurrency int) {
	wg := sync.WaitGroup{}
	ch := make(chan email.EmailInput)
	var validate = validator.New()

	// Start workers
	for i := 0; i < concurrency; i++ {
		go func() {
			for e := range ch {
				// Validate email
				if err := validate.Struct(e); err != nil {
					log.Printf("[VALIDATION ERROR] Skipping email to %v: %v", e.T, err)
					wg.Done()
					continue // skip sending invalid email
				}

				// Send email
				if err := messenger.SendEmail(e); err != nil {
					log.Printf("[ERROR] Sending to %v: %v", e.T, err)
				}
				wg.Done()
			}
		}()
	}

	// Feed emails
	for _, e := range req.Emails {
		wg.Add(1)
		ch <- e
	}

	wg.Wait()
	close(ch)
}
