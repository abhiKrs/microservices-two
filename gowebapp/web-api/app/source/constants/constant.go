package constants

// Enum SourceType
type SourceType uint8

func (d SourceType) String() string {
	strings := [...]string{
		"magicLink",
		"password",
		"google",
		"docker",
		"nginx",
		"dokku",
		"fly.io",
		"heroku",
		"ubuntu",
		"vercel",
		".net",
		"apache2",
		"cloudflare",
		"java",
		"python",
		"php",
		"postgresql",
		"redis",
		"ruby",
		"mongodb",
		"mysql",
		"http",
		"vector",
		"fluentbit",
		"fluentd",
		"logstash",
		"rsyslog",
		"render",
		"syslog-ng",
		"demo",
	}

	if d < Kubernetes || d > Demo {
		return "Unknown"
	}
	return strings[d-1]
}

func (d SourceType) EnumIndex() int {
	return int(d)
}
