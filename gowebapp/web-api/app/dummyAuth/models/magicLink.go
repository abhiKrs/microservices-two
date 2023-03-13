package models

import (
	"encoding/json"
	"io"
	"time"
)

type MagicLinkRequest struct {
	Email string `json:"email" validate:"email,required"`
}

type MagicLinkResponse struct {
	IsSuccessful bool      `json:"isSuccessful"`
	Message      []string  `json:"message,omitempty"`
	Email        string    `json:"email"`
	ExpiryTime   time.Time `json:"expiryTime"`
}

func (r *MagicLinkRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
