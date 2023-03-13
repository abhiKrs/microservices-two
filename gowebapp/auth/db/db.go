package db

import (
	"context"
	// "net/http"
	"fmt"
	"log"
	"os"

	// "log"
	// "auth/config"
	// "auth/db/schema"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	sqlBindingName string = "postgresdb"
)

func AddUuid() error {

	fmt.Println("Processing db...")

	client, err := dapr.NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	a := "uuid-ossp"
	query := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS '%v';", a)
	fmt.Println(query)

	// Insert query using Dapr output binding via Dapr SDK
	in := &dapr.InvokeBindingRequest{
		Name:      sqlBindingName,
		Operation: "exec",
		Data:      []byte(""),
		Metadata:  map[string]string{"sql": query},
	}

	err = client.InvokeOutputBinding(ctx, in)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("? Executedted Successfully to the Database")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return nil

}

func DBConnectUrl() string {
	// var err error

	var daprHost, daprHttpPort string
	var okHost, okPort bool

	if daprHost, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		daprHost = "http://localhost"
	}

	if daprHttpPort, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		daprHttpPort = "3503"
	}

	var daprUrl string = daprHost + ":" + daprHttpPort + "/v1.0/bindings/" + sqlBindingName
	fmt.Printf("Connect the Database @ %v", daprUrl)

	return daprUrl
}
