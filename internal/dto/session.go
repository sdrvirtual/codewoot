package dto

type CreateSession struct {
	SessionID   *string `json:"session_id"`
	Description *string `json:"description"`
	Chatwoot    struct {
		InboxID   int    `json:"inbox_id"`
		AccountID int    `json:"account_id"`
		Token     string `json:"token"`
	} `json:"chatwoot"`
}
