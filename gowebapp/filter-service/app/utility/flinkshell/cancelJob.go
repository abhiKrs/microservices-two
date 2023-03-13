package flinkshell

import (
	log "filter-service/app/utility/logger"
	"os/exec"
)

func CancelJob(oldJob string) error {
	cmd := exec.Command("./bin/flink", "cancel", oldJob)
	stdout, err := cmd.Output()

	if err != nil {
		log.ErrorLogger.Println(err.Error())
		// return
		return err
	}

	log.DebugLogger.Println(stdout)
	return nil
}
