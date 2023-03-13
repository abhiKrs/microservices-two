package model

import "time"

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

type Message struct {
	Type int  `json:"type,omitempty"`
	Data Data `json:"data,omitempty"`
}
