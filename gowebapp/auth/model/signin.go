package model

import (
	"encoding/json"
	"io"
)

type SignInRequest struct {
	Email      *string `json:"email"`
	Credential string  `json:"credential" validate:"required"`
	AuthType   uint8   `json:"authType" validate:"required"`
}

type SignInResponse struct {
	IsSuccessful bool        `json:"isSuccssful"`
	Code         int         `json:"code"`
	Email        string      `json:"email,omitempty"`
	UserBody     ProfileBody `json:"userBody,omitempty"`
	BearerToken  string      `json:"bearerToken,omitempty"`
	Message      []string    `json:"message,omitempty"`
}

func (r *SignInRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
