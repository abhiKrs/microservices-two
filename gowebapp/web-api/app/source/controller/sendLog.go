package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	// "web-api/app/source/models"
	"web-api/app/source/constants"
	"web-api/app/source/models"
	myKafka "web-api/app/utility/kafka"
	logs "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	// "github.com/google/uuid"
)

func (sc *SourceController) SendLog(w http.ResponseWriter, r *http.Request) {
	var err error

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		logs.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(bearerToken, "bearer_")
	sourceToken := strings.Join(splitToken[1:], "")

	// detoken := strings.Split(sourceToken, "_")

	// var sourceId, profileId uuid.UUID

	// if len(detoken) != 2 {
	// 	logs.ErrorLogger.Println("wrong token")
	// 	respond.Error(w, http.StatusBadRequest, errors.New("wrong token"))
	// 	return
	// } else {
	// 	sourceId, err = uuid.Parse(detoken[0])
	// 	if err != nil {
	// 		logs.ErrorLogger.Println(err)
	// 		respond.Errors(w, http.StatusBadRequest, []string{"wrong token", err.Error()})
	// 		return
	// 	}
	// 	profileId, err = uuid.Parse(detoken[1])
	// 	if err != nil {
	// 		logs.ErrorLogger.Println(err)
	// 		respond.Errors(w, http.StatusBadRequest, []string{"wrong token", err.Error()})
	// 		return
	// 	}
	// }

	decoded, err := base64.StdEncoding.DecodeString(sourceToken)
	if err != nil {
		logs.ErrorLogger.Println(err)
		respond.Errors(w, http.StatusUnauthorized, []string{"wrong token", err.Error()})
		return
	}
	logs.DebugLogger.Println(string(decoded))
	str := string(decoded)
	s := strings.Split(str, ",")

	sourceKV := strings.Split(s[0], ":")
	logs.DebugLogger.Printf("source id passed: %s", sourceKV[1])
	sourceId := sourceKV[1]

	typeKV := strings.Split(s[1], ":")
	logs.DebugLogger.Printf("source type passed: %s", typeKV[1])
	n, err := strconv.ParseUint(typeKV[1], 10, 8)
	if err != nil {
		logs.ErrorLogger.Println(err)
		respond.Errors(w, http.StatusUnauthorized, []string{"wrong token", err.Error()})
		return
	}

	sourceType := constants.SourceType(n)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logs.ErrorLogger.Printf("could not read body: %s\n", err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	// var getSource *models.SingleSourceResponse
	// getSource, err = sc.sourceService.GetSourceById(sourceId, profileId)
	// if err != nil {
	// 	logs.ErrorLogger.Println(err)
	// 	respond.Errors(w, http.StatusBadRequest, []string{"wrong token", err.Error()})
	// 	return
	// }

	sourceTopic := fmt.Sprintf("source_topic_%s", sourceId)

	rh, err := myKafka.NewRecordHandler(sourceTopic)
	if err != nil {
		logs.ErrorLogger.Printf("Failed to create serializer: %s\n", err)
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	var record interface{}

	if sourceType == constants.Docker {
		record = models.DockerVectorDefaultLogs{}
		err = json.Unmarshal(body, &record)
		if err != nil {
			logs.WarningLogger.Printf("could not unmarshal body: %s\n", err)
			respond.Errors(w, http.StatusBadRequest, []string{err.Error()})
			return
		}

	} else {
		logs.ErrorLogger.Printf("source type %v: not supported yet", sourceType.String())
		respond.Error(w, http.StatusInternalServerError, fmt.Errorf("source type {%v:%v} not supported yet", sourceType.String(), sourceType.EnumIndex()))
		return
	}

	logs.DebugLogger.Println(record)

	errs := validate.Validate(sc.validator, record)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	payload, err := rh.Ser.Serialize(rh.Topic, &record)
	if err != nil {
		logs.WarningLogger.Printf("Failed to serialize payload: %s\n", err)
		respond.Errors(w, http.StatusBadRequest, []string{err.Error()})
		return
	}

	deliveryChan := make(chan kafka.Event)

	// TODO Patition
	err = rh.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &rh.Topic, Partition: kafka.PartitionAny},
			Value:          payload,
		},
		deliveryChan,
	)

	if err != nil {
		logs.WarningLogger.Printf("Produce failed: %v\n", err)
		respond.Errors(w, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		logs.ErrorLogger.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
		respond.Errors(w, http.StatusInternalServerError, []string{err.Error()})
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic,
			m.TopicPartition.Partition,
			m.TopicPartition.Offset,
		)
		respond.Json(w, http.StatusAccepted, map[string]bool{"isSuccessful": true})
	}

	close(deliveryChan)

	// var records []DockerVectorDefaultLogs

	// err = json.Unmarshal(body, &records)
	// if err != nil {
	// 	logs.WarningLogger.Printf("could not unmarshal body: %s\n", err)
	// 	respond.Errors(w, http.StatusBadRequest, []string{err.Error()})
	// 	return
	// }

	// payload, err := ser.Serialize(topic, &value)
	// if err != nil {
	// 	fmt.Printf("Failed to serialize payload: %s\n", err)
	// 	os.Exit(1)
	// }

	// for _, record := range records {
	// 	payload, err := rh.Ser.Serialize(rh.Topic, &record)
	// 	if err != nil {
	// 		logs.WarningLogger.Printf("Failed to serialize payload: %s\n", err)
	// 		continue
	// 	}

	// 	err = rh.Producer.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &rh.Topic, Partition: 0},
	// 		Value:          payload,
	// 	}, nil)
	// 	if err != nil {
	// 		logs.WarningLogger.Printf("Produce failed: %v\n", err)
	// 	}
	// }

	// _, err = io.WriteString(w, "Success!\n")
	// if err != nil {
	// 	logs.WarningLogger.Printf("could not write: %s\n", err)
	// 	return
	// }

}

// func (pc *SourceController) SendLog(w http.ResponseWriter, r *http.Request) {

// 	// var reqBody models.SourceLog
// 	// err := reqBody.Bind(r.Body)
// 	// if err != nil {
// 	// 	logs.ErrorLogger.Println(err)
// 	// 	respond.Error(w, http.StatusBadRequest, err)
// 	// 	return
// 	// }

// 	// logs.DebugLogger.Println(reqBody)

// 	// errs := validate.Validate(pc.validator, reqBody)
// 	// if errs != nil {
// 	// 	respond.Errors(w, http.StatusBadRequest, errs)
// 	// 	return
// 	// }

// 	// if sourceId == "" {
// 	// 	response := models.SourceLogResponse{IsSuccessful: false, Message: []string{"Wrong or missing profile Id"}}
// 	// 	respond.Json(w, http.StatusBadRequest, response)
// 	// 	return
// 	// }

// 	var reqBody interface{}
// 	json.NewDecoder(r.Body).Decode(&reqBody)

// 	bearerToken := r.Header.Get("Authorization")
// 	if bearerToken == "" {
// 		logs.ErrorLogger.Println("Missing Authorization Header")
// 		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
// 	}

// 	// TODO use JWT token

// 	splitToken := strings.Split(bearerToken, "bearer_")
// 	sourceToken := strings.Join(splitToken[1:], "")

// 	// source topic
// 	// sourceTopic := "source_topic_" + sourceToken
// 	// userId for now
// 	sourceTopic := "user_topic_" + sourceToken

// 	err := kafka.SendLog(sourceTopic, reqBody)
// 	if err != nil {
// 		respond.Error(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	respond.Json(w, http.StatusAccepted, map[string]bool{"isSuccessful": true})

// }
