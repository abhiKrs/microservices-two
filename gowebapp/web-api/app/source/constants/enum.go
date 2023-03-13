package constants

// Declare related constants for each SourceType starting with index 0
const (
	Kubernetes SourceType = iota + 1 // EnumIndex = 1
	AWS                              // EnumIndex = 2
	JavaScript                       // EnumIndex = 3
	Docker
	Nginx
	Dokku
	FlyDotio
	Heroku
	Ubuntu
	Vercel
	DotNET
	Apache2
	Cloudflare
	Java
	Python
	PHP
	PostgreSQL
	Redis
	Ruby
	MongoDB
	MySQL
	HTTP
	Vector
	FluentBit
	Fluentd
	Logstash
	RSyslog
	Render
	SyslogNg
	Demo
)
