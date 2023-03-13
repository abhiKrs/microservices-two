package flinkshell

import (
	"bytes"
	"encoding/json"

	// "errors"
	"filter-service/app/models"
	log "filter-service/app/utility/logger"
	"fmt"
	"os"
	"os/exec"
	// "strings"
)

type FlinkJobArgs struct {
	Levels          *[]string `json:"levels,omitempty"`
	Platforms       *[]string `json:"platforms,omitempty"`
	StartTime       *string   `json:"start_time,omitempty"`
	EndTime         *string   `json:"end_time,omitempty"`
	MessageContains *string   `json:"message_contains,omitempty"`
	Sinktopic       string    `json:"sinktopic"`
	Sourcetopic     string    `json:"sourcetopic"`
}

func CreateJob(
	streamFilter models.StreamFilter,
	sourceTopic string,
	customTopic string,
) (string, error) {
	log.DebugLogger.Println("started new job creation")
	f, err := os.Create(os.TempDir() + "/tmpfile-filtered-job-aaplication.txt")

	if err != nil {
		log.DebugLogger.Println(err)
		return "", err
	}

	// close and remove the temporary file at the end of the program
	defer f.Close()
	// defer os.Remove(f.Name())

	// write data to the temporary file
	flinkJob := FlinkJobArgs{
		Levels:    streamFilter.Level,
		Platforms: streamFilter.Platform,
		// StartTime:       startTime,
		// EndTime:         endTime,
		MessageContains: streamFilter.Text,
		Sourcetopic:     sourceTopic,
		Sinktopic:       customTopic,
	}
	var startTime string
	var endTime string

	if streamFilter.Date != nil {
		if streamFilter.Date.StartDate != nil {
			date := streamFilter.Date.StartDate
			// log.DebugLogger.Println(date)
			startTime = date.Format("2006-01-02 15:04:05")
			// startTime = date.Format("25-02-2023 09:54:24")
			log.DebugLogger.Println(startTime)
			log.DebugLogger.Println(&startTime)
			flinkJob.StartTime = &startTime
		}
		if streamFilter.Date.EndDate != nil {
			date := streamFilter.Date.EndDate
			log.DebugLogger.Println(date)
			endTime = date.Format("2006-01-02 15:04:05")
			// endTime = date.Format("25-02-2023 09:54:24")
			flinkJob.EndTime = &endTime
		}
	}

	// flinkJob := FlinkJobArgs{
	// 	Levels:          streamFilter.Level,
	// 	Platforms:       streamFilter.Platform,
	// 	StartTime:       startTime,
	// 	EndTime:         endTime,
	// 	MessageContains: streamFilter.Text,
	// 	Sourcetopic:     sourceTopic,
	// 	Sinktopic:       customTopic,
	// }

	flinkByteData, err := json.Marshal(flinkJob)
	log.DebugLogger.Println(string(flinkByteData))
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", err
	}

	// data := []byte("abc abc abc")
	if _, err := f.Write(flinkByteData); err != nil {
		log.DebugLogger.Println(err)
		return "", err
	}
	log.InfoLogger.Println(f.Name())

	cmd := exec.Command(
		"./bin/flink", "run",
		// "--target", "kubernetes-session",
		// "--Dkubernetes.namespace", "logfire-local",
		// "--hostname", "flink-jobmanager", "--port", "8081",
		// "-Dkubernetes.cluster-id=my-first-flink-cluster",
		"--jobmanager", "flink-jobmanager:8081",
		"--python", "main.py",
		"--pyFiles", "file:///"+f.Name()+",file:///opt/flink/bin/flink-sql-connector-kafka-1.16.1.jar",
		"--jarfile", "bin/flink-sql-connector-kafka-1.16.1.jar",
	)
	go func(cmd *exec.Cmd) {
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			log.DebugLogger.Println(fmt.Sprint(err) + ": " + stderr.String())
			// return "", err
			// return
		}
		log.DebugLogger.Println(out.String())
		// res := strings.Split(out.String(), " ")
		// sliceb := bytes.Split(out.Bytes(), []byte(" "))

		i := bytes.Index(out.Bytes(), []byte("JobID"))
		if i != -1 {
			log.DebugLogger.Println(string(out.Bytes()[i]))
			s := bytes.SplitAfter(out.Bytes(), []byte("JobID"))
			a := s[len(s)-1]
			log.DebugLogger.Println(string(a[1:]))
			// b := bytes.Split(a, []byte(" "))
			// log.DebugLogger.Println(string(b[1]))
			// return string(b[1]), nil
			// return string(a[1:]), nil
		} else {
			log.ErrorLogger.Println("no job created")
			// return "", errors.New("JobID not found")
		}
		// log.DebugLogger.Println(out.)

	}(cmd)

	return "jobId", nil
}
