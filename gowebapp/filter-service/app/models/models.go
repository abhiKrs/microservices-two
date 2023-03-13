package models

import (
	"encoding/json"
	"io"
	"time"
)

type Response struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
	StreamId     string   `json:"streamId,omitempty"`
}

type ErrorsMessage struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message"`
}

type StreamFilterResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
	StreamId     string   `json:"streamId,omitempty"`
}

type DateInterval struct {
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
}

type StreamFilter struct {
	ProfileId string        `json:"profileId"`
	TeamId    string        `json:"teamId"`
	SourceId  string        `json:"sourceId"`
	Level     *[]string     `json:"level,omitempty"`
	Date      *DateInterval `json:"date,omitempty"`
	Platform  *[]string     `json:"platform,omitempty"`
	Pid       string        `json:"pid,omitempty"`
	Text      *string       `json:"text,omitempty"`
}

func (f *StreamFilter) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(f)
}
