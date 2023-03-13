package main

import (
	// "encoding/json"

	"errors"
	"fmt"

	// "fmt"
	"net/http"
	"strings"

	"filter-service/app/models"
	// "filter-service/app/utility/flinkshell"
	"filter-service/app/utility/kafka"
	logs "filter-service/app/utility/logger"

	myRedis "filter-service/app/utility/redis"
	"filter-service/app/utility/respond"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func StreamFilterHandler(w http.ResponseWriter, r *http.Request) {

	logs.DebugLogger.Println("inside filter endpoint")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		logs.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusBadRequest, errors.New("missing authorization header"))
		return
	}

	// TODO use JWT token
	splitToken := strings.Split(reqToken, "bearer_")

	profileId := strings.Join(splitToken[1:], "")
	logs.InfoLogger.Println(profileId)

	if profileId == "" {
		response := models.StreamFilterResponse{IsSuccessful: false, Message: []string{"Wrong or missing profileId in query parameter"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	pId, err := uuid.Parse(profileId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.StreamFilterResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// payload, err := pc.sourceService.GetSourcesByProfileId(pId)
	// if err != nil {
	// 	logs.ErrorLogger.Println(err)
	// 	response := models.StreamFilterResponse{IsSuccessful: false, Message: []string{err.Error()}}
	// 	Json(w, http.StatusBadRequest, response)
	// 	return
	// }
	// Json(w, http.StatusAccepted, payload)
	// return

	logs.DebugLogger.Println("profileId: ", pId)
	var req models.StreamFilter
	// err := json.NewDecoder(r.Body).Decode(r)
	err = req.Bind(r.Body)
	if err != nil {
		logs.ErrorLogger.Println(err)
		res := models.Response{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusInternalServerError, res)
		return
	}

	// logs.DebugLogger.Printf("%v, %v, %v, %v", *req.Date, *req.Level, *req.Platform, *req.Text)
	logs.DebugLogger.Printf("%v, %v, %v, %v", req.Date, req.Level, req.Platform, req.Text)
	// logs.DebugLogger.Println(req.Date.EndDate.Format("2023-02-15 07:16:29.097"))

	// levelFilter := req.Level
	// sourceFilter := req.SourcesFilter

	// generate random topic name
	streamId := uuid.New()
	logs.DebugLogger.Println(&req)

	filterTopic := "filter_topic_" + streamId.String()
	logs.DebugLogger.Println("created topic named: ", filterTopic)

	// Crete kafka topic for filter
	err = kafka.CreateTopic(filterTopic)
	if err != nil {
		logs.ErrorLogger.Println(err)
		res := models.Response{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusInternalServerError, res)
		return
	}

	// Create redis store custome filter per user
	logs.DebugLogger.Println("topic creation success")
	err = myRedis.AddFilterPerUser(filterTopic, req)
	if err != nil {
		logs.DebugLogger.Println(err)
		err = kafka.DeleteTopic(filterTopic)
		if err != nil {
			logs.WarningLogger.Println(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error: %v", err.Error())))
		return
	}

	// get jobId from redis store
	// oldJobId := myRedis.GetJobIdByUser(profileId)
	// logs.DebugLogger.Println(oldJobId)
	// if oldJobId != "" {
	// 	logs.DebugLogger.Println(oldJobId)
	// 	// kill old job
	// 	err := flinkshell.CancelJob(oldJobId)
	// 	if err != nil {
	// 		logs.DebugLogger.Println(err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(fmt.Sprintf("error: %v", err.Error())))
	// 		return
	// 	}
	// }

	// ------------------------------------------------------
	// logs.DebugLogger.Println("next to run flinkshell")
	// // Submit flink-job
	// // get jobId and store to redis
	// sourceTopic := "user_topic_" + profileId

	// jobId, err := flinkshell.CreateJob(req, sourceTopic, filterTopic)
	// if err != nil {
	// 	logs.ErrorLogger.Println(err)
	// 	res := models.Response{IsSuccessful: false, Message: []string{err.Error()}}
	// 	respond.Json(w, http.StatusInternalServerError, res)
	// 	return
	// }
	// logs.DebugLogger.Println("got jobId: ", jobId)

	// --------------------------------------------------
	// if err != nil {
	// 	logs.ErrorLogger.Println(err)
	// }
	// if jobId == "" {
	// 	logs.DebugLogger.Println("no job id returned")
	// 	// w.WriteHeader(http.StatusInternalServerError)
	// 	// w.Write([]byte(fmt.Sprintf("error: %v", err.Error())))
	// 	// return
	// }

	// logs.DebugLogger.Println("done running flinkshell")
	// err = myRedis.AddUserJobIdPair(jobId, "user_"+profileId)
	// if err != nil {
	// 	logs.DebugLogger.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(fmt.Sprintf("error: %v", err.Error())))
	// 	return
	// }
	// logs.InfoLogger.Println(jobId)
	res := models.Response{IsSuccessful: true, StreamId: streamId.String()}
	respond.Json(w, http.StatusAccepted, res)
	// logs.InfoLogger.Println(<-sol)
	// logs.InfoLogger.Println(<-errCh)

	// return

}

func main() {
	port := ":8005"

	logs.DebugLogger.Println("Started Go Flink filter")
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "http://localhost"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello From Filter-service!"))
	})

	router.Post("/", StreamFilterHandler)

	logs.InfoLogger.Printf("Listening at port: %v", port)
	if err := http.ListenAndServe(port, router); err != nil {
		logs.ErrorLogger.Println(err)
		// logs.ErrorLogger.Fatal(err)
	}

}
