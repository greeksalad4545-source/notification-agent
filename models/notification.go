package models

type Notification struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Channel        string `json:"channel"`
	Subject        string `json:"subject"`
	Message        string `json:"message"`
	RecipientEmail string `json:"recipient_email"`
	RecipientPhone string `json:"recipient_phone"`
}
