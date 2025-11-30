package dto

type WebhookPayload struct {
	Status   string `json:"status"`
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url,omitempty"`
	Code     string `json:"code,omitempty"`
	Error    string `json:"error,omitempty"`
}
