package events

type EmailMessage struct {
	Type     string `json:"type"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Retry    int    `json:"retry"`
	MaxRetry int    `json:"max_retry"`
}
