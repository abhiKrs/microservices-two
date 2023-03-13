package flinkshell

import (
	"os/exec"
	log "web-api/app/utility/logger"
)

func CreateJob() string {
	// ./bin/flink run --target kubernetes-session --python filter.py --jarfile connector.jar
	// ./bin/flink run \
	// --jobmanager <jobmanagerHost>:8081 \
	// --python examples/python/table/word_count.py

	cmd := exec.Command("./bin/flink", "run", "--python", "filter.py", "--jarfile", "connector.jar")
	stdout, err := cmd.Output()

	if err != nil {
		log.ErrorLogger.Println(err.Error())
		// return
		return ""
	}

	log.DebugLogger.Println(stdout)
	return "getJobId"
}
