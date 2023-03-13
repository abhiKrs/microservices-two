package models

import (
	"encoding/json"
	"io"

	"source/app/api/schema"
)

type CreateSourceRequest struct {
	Name         string `json:"name" validate:"required"`
	PlatformType uint8  `json:"platformType" validate:"required"`
}

type CreateSourceResponse struct {
	IsSuccessful bool                `json:"isSuccessful"`
	Message      []string            `json:"message,omitempty"`
	SourceBody   ProfileBodyResponse `json:"sourceBody,omitempty"`
}

type ProfileBodyResponse struct {
	schema.Source
}

func (r *CreateSourceRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
