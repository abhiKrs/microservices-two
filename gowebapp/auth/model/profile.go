package model

import (
	"auth/db/schema"
	"encoding/json"
	"io"

	"github.com/google/uuid"
)

type ProfileBody struct {
	ProfileId      string     `json:"profileId,omitempty"`
	AccountId      *uuid.UUID `json:"accountId,omitempty"`
	Onboarded      *bool      `json:"onboarded,omitempty"`
	Email          string     `json:"email,omitempty"`
	AccessApproved *bool      `json:"accessApproved,omitempty"`
	schema.ProfileBase
}

type ProfileUpdateRequest struct {
	schema.ProfileBase
	Password *string `json:"password,omitempty"`
}

func (r *ProfileUpdateRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

type ProfileResponse struct {
	IsSuccessful bool                `json:"isSuccssful"`
	UserBody     ProfileBodyResponse `json:"userBody,omitempty"`
	Message      []string            `json:"message,omitempty"`
}

type ProfileBodyResponse struct {
	schema.Profile
	Email          *string `json:"email,omitempty"`
	AccessApproved bool    `json:"accessApproved"`
}
