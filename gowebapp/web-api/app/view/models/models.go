package models

import (
	"encoding/json"
	"io"
	"time"
	// "github.com/google/uuid"
)

type DateInterval struct {
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

type CreateViewRequest struct {
	TeamId        string        `json:"teamId" validate:"required"`
	StreamId      string        `json:"streamId,omitempty"`
	SourcesFilter *[]string     `json:"sourcesFilter,omitempty"`
	LevelFilter   *[]string     `json:"levelFilter,omitempty"`
	DateFilter    *DateInterval `json:"dateFilter,omitempty"`
	Name          string        `json:"name" validate:"required"`
}

type UpdateViewRequest struct {
	Name          *string       `json:"name,omitempty"`
	StreamId      *string       `json:"streamId,omitempty"`
	SourcesFilter *[]string     `json:"sourcesFilter,omitempty"`
	LevelFilter   *[]string     `json:"levelFilter,omitempty"`
	DateFilter    *DateInterval `json:"dateFilter,omitempty"`
}

func (r *UpdateViewRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

type UpdateViewResponse struct {
	IsSuccessful bool             `json:"isSuccessful"`
	Message      []string         `json:"message,omitempty"`
	Data         ViewBodyResponse `json:"data,omitempty"`
}

type CreateViewResponse struct {
	IsSuccessful bool             `json:"isSuccessful"`
	Message      []string         `json:"message,omitempty"`
	Data         ViewBodyResponse `json:"data,omitempty"`
}

type DeleteViewResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

type GetViewResponse struct {
	IsSuccessful bool             `json:"isSuccessful"`
	Message      []string         `json:"message,omitempty"`
	Data         ViewBodyResponse `json:"data,omitempty"`
}

type GetAllViewResponse struct {
	IsSuccessful bool               `json:"isSuccessful"`
	Message      []string           `json:"message,omitempty"`
	Data         []ViewBodyResponse `json:"data,omitempty"`
}

type ViewBodyResponse struct {
	CreateViewRequest
	CreateDate time.Time `json:"createDate"`
}

func (r *CreateViewRequest) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
