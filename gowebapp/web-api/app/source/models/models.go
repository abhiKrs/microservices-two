package models

import (
	"encoding/json"
	"io"

	// "time"

	"web-api/app/source/schema"
	// "github.com/google/uuid"
)

type CreateSourceRequest struct {
	Name       string `json:"name" validate:"required"`
	SourceType uint8  `json:"sourceType" validate:"required"`
	TeamId     string `json:"teamId" validate:"required"`
}

func (r *CreateSourceRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

type UpdateSourceRequest struct {
	Name string `json:"name" validate:"required"`
}

func (r *UpdateSourceRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

type SingleSourceResponse struct {
	IsSuccessful bool       `json:"isSuccessful"`
	Message      []string   `json:"message,omitempty"`
	Data         SourceBody `json:"data,omitempty"`
	// Data         schema.Source `json:"data,omitempty"`
}

type AllSourceResponse struct {
	IsSuccessful bool         `json:"isSuccessful"`
	Message      []string     `json:"message,omitempty"`
	Data         []SourceBody `json:"data,omitempty"`
	// Data         []schema.Source `json:"data,omitempty"`
}

type SourceBody struct {
	schema.Source
	SourceToken string `json:"sourceToken,omitempty"`
}
