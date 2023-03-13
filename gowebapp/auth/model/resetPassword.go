package model

import (
	"encoding/json"
	"io"
	// "time"
)

type ResetPasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

type ResetPasswordResponse struct {
	IsSuccessful bool   `json:"isSuccessful"`
	UserId       string `json:"userId,omitempty"`
	BearerToken  string `json:"bearerToken,omitempty"`
	Message      string `json:"message,omitempty"`
}

func (r *ResetPasswordRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
