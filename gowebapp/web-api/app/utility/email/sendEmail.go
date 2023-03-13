package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	constants "web-api/app/constants"
	dbConstants "web-api/app/dummyAuth/constants"
	log "web-api/app/utility/logger"
)

type EmailRequestBody struct {
	Email       string                       `json:"email"`
	Token       *string                      `json:"token,omitempty"`
	EmailType   constants.EmailType          `json:"emailType"`
	LinkType    *dbConstants.DBMagicLinkType `json:"linkType,omitempty"`
	SuccessType *constants.SuccessType       `json:"successType"`
}

type EmailResponseBody struct {
	IsSuccessful bool `json:"isSuccessful"`
}

func SendEmail(email string, token *string, messageType constants.EmailType, linkType *dbConstants.DBMagicLinkType, successType *constants.SuccessType) error {
	bd := EmailRequestBody{
		Email:       email,
		Token:       token,
		EmailType:   messageType,
		LinkType:    linkType,
		SuccessType: successType,
	}
	var NOTIFICATION_HOST, NOTIFICATION_HTTP_PORT string
	var okHost, okPort bool
	if NOTIFICATION_HOST, okHost = os.LookupEnv("NOTIFICATION_HOST"); !okHost {
		NOTIFICATION_HOST = "http://notification"
	}
	if NOTIFICATION_HTTP_PORT, okPort = os.LookupEnv("NOTIFICATION_HTTP_PORT"); !okPort {
		NOTIFICATION_HTTP_PORT = "80"
	}

	// order := "{\"email\":" + strconv.Itoa(i) + "}"
	bdJSON, err := json.Marshal(bd)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", NOTIFICATION_HOST+":"+NOTIFICATION_HTTP_PORT+"/email", bytes.NewBuffer(bdJSON))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Adding app id as part of th header
	// req.Header.Add("dapr-app-id", "notification")

	// Invoking a service
	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	var resBody EmailResponseBody
	err = json.Unmarshal(result, &resBody)
	if err != nil {
		log.DebugLogger.Printf("------{%v}-----", err)
		return err
	}

	fmt.Println("Response from server:", resBody)
	return nil

}

func SendEmailviaDapr(email string, token *string, messageType constants.EmailType, linkType *dbConstants.DBMagicLinkType, successType *constants.SuccessType) error {
	bd := EmailRequestBody{
		Email:       email,
		Token:       token,
		EmailType:   messageType,
		LinkType:    linkType,
		SuccessType: successType,
	}
	var DAPR_HOST, DAPR_HTTP_PORT string
	var okHost, okPort bool
	if DAPR_HOST, okHost = os.LookupEnv("DAPR_HOST"); !okHost {
		DAPR_HOST = "http://localhost"
	}
	if DAPR_HTTP_PORT, okPort = os.LookupEnv("DAPR_HTTP_PORT"); !okPort {
		DAPR_HTTP_PORT = "3500"
	}

	// order := "{\"email\":" + strconv.Itoa(i) + "}"
	bdJSON, err := json.Marshal(bd)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", DAPR_HOST+":"+DAPR_HTTP_PORT+"/email", bytes.NewBuffer(bdJSON))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Adding app id as part of th header
	req.Header.Add("dapr-app-id", "notification")

	// Invoking a service
	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	var resBody EmailResponseBody
	err = json.Unmarshal(result, &resBody)
	if err != nil {
		log.DebugLogger.Printf("------{%v}-----", err)
		return err
	}

	fmt.Println("Response from server:", resBody)
	return nil

}
