package models

import (
	"encoding/json"
	"io"
	"time"
)

type Runtime struct {
	Function string `json:"function,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	ThreadId int    `json:"threadId,omitempty"`
}

type System struct {
	Pid         int    `json:"pid,omitempty"`
	ProcessName string `json:"processName,omitempty"`
}

type Context struct {
	Runtime Runtime `json:"runtime,omitempty"`
	System  System  `json:"system,omitempty"`
}

type Data struct {
	Ts time.Time `json:"ts,omitempty"`
	// LogLevel  int     `json:"logLevel,omitempty"`
	LogString string  `json:"logString,omitempty"`
	Message   string  `json:"message,omitempty"`
	Context   Context `json:"context,omitempty"`
	Source    string  `json:"source,omitempty"`
}

type SourceLog struct {
	Type int    `json:"type,omitempty"`
	Data []Data `json:"data" validate:"required"`
}

func (r *SourceLog) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}

type SourceLogResponse struct {
	IsSuccessful bool     `json:"isSuccessful"`
	Message      []string `json:"message,omitempty"`
}

type DockerVectorDefaultLogs struct {
	SourceType string `json:"source_type,omitempty"`
	Level      string `json:"level,omitempty"`
	Message    string `json:"message,omitempty"`
	// Timestamp  string `json:"timestamp,omitempty"`
	Timestamp string `json:"dt,omitempty"`

	ContainerCreatedAt string            `json:"container_created_at,omitempty"`
	ContainerId        string            `json:"container_id,omitempty"`
	ContainerName      string            `json:"container_name,omitempty"`
	Host               string            `json:"host,omitempty"`
	Stream             string            `json:"stream,omitempty"`
	Label              map[string]string `json:"label,omitempty"`
}

func (r *DockerVectorDefaultLogs) Bind(body io.ReadCloser) error {
	return json.NewDecoder(body).Decode(r)
}
