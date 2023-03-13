package utility

import (
	"math/rand"
	"strings"
	"time"

	"stream_source/app/logs/model"
)

func getFn() []string {
	return []string{"my_service", "my_controller", "my_dao", "my_server", "my_fun"}
}

func getLvl() []string {
	return []string{"debug", "info", "warning", "error", "fatal", "log", "statement", "trace"}
}

func getPName() []string {
	return []string{"Python", "MainProcess", "Docker", "Kafka"}
}

func getSName() []string {
	return []string{"Docker", "AWS", "GCP", "Server", "Kubernetes"}
}

func GenerateData(msg string) model.Data {
	timeStamp := time.Now()
	fun := getFn()
	lvl := getLvl()
	pName := getPName()
	sName := getSName()
	rand.Seed(time.Now().UnixNano())

	var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	payload := model.Data{
		Ts:        timeStamp,
		Message:   msg,
		LogString: lvl[rand.Intn(len(lvl)-1)],
		Source:    sName[rand.Intn(len(sName)-1)],
		Context: model.Context{
			Runtime: model.Runtime{
				Function: fun[rand.Intn(len(fun)-1)],
				File:     randomString(10, alphabet),
				Line:     rand.Intn(50) + 1,
				ThreadId: rand.Intn(5000) + 1,
			},
			System: model.System{
				Pid:         rand.Intn(15000) + 1,
				ProcessName: pName[rand.Intn(len(pName)-1)],
			},
		},
	}
	// fmt.Printf("This is pretty random: %v\n", payload)
	return payload
}

func randomString(n int, alphabet []rune) string {

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}
