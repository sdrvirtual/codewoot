package dto

import "net/http"

type APIErrorResponse struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

func NewAPIErrorResponse(message string, description string) *APIErrorResponse {
	return &APIErrorResponse{
		Message:     message,
		Description: description,
	}
}

func (rd *APIErrorResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }
