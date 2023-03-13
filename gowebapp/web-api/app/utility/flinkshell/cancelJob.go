package flinkshell

import (
	"os/exec"
	log "web-api/app/utility/logger"
)

func CancelJob(oldJob string) error {
	cmd := exec.Command("./bin/flink", "cancel", "oldJobId")
	stdout, err := cmd.Output()

	if err != nil {
		log.ErrorLogger.Println(err.Error())
		// return
		return err
	}

	log.DebugLogger.Println(stdout)
	return nil
}
