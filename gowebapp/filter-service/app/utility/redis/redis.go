package controller

import (
	"encoding/json"
	"filter-service/app/models"
	log "filter-service/app/utility/logger"

	"github.com/go-redis/redis"
)

type TeamSource struct {
	TeamId   string `json:"teamId"`
	SourceId string `json:"sourceId"`
}

// Create Redis Client

func MyRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-service.redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

var client = MyRedisClient()

// --------{Source : Team}--------------

func AddSourceTeamPair(sourceTopic string, teamTopic string) error {

	err := client.Set(sourceTopic, teamTopic, 0).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func GetTeamBySource(sourceTopic string) string {

	res, err := client.Get(sourceTopic).Result()
	if err != nil {
		log.ErrorLogger.Println(err)
		return ""
	}
	return res
}

// --------{Team : [Source]}--------------

func PushSourceToTeamSet(sourceTopic string, teamTopic string) error {

	// res, err := client.Get(teamTopic).Result()
	// if err == redis.Nil {
	// 	err := client.SAdd(teamTopic, sourceTopic).Err()
	// 	if err != nil {
	// 		log.ErrorLogger.Println(err)
	// 		return err
	// 	}
	// 	return nil
	// }

	// if err != nil {
	// 	log.ErrorLogger.Println(err)
	// 	return err
	// }

	err := client.SAdd(teamTopic, sourceTopic).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func GetSourcesByTeam(teamTopic string) string {

	res, err := client.Get(teamTopic).Result()
	if err != nil {
		log.ErrorLogger.Println(err)
		return ""
	}
	return res
}

// --------{User : Filter}--------------

type MyFilter struct {
	StreamTopic string   `json:"streamTopic,omitempty"`
	LevelList   []string `json:"levelList,omitempty"`
	SourceList  []string `json:"sourceList,omitempty"`
}

func getUserFilterKey(profileId string) string {
	return "user_filter_" + profileId
}

func AddFilterPerUser(customTopic string, filterReq models.StreamFilter) error {

	val, err := client.Del(getUserFilterKey(filterReq.ProfileId)).Result()
	log.DebugLogger.Println(val)
	if err == redis.Nil {
		log.DebugLogger.Println("No previous active filter")
		// return err
	}

	if err != nil {
		log.ErrorLogger.Println(err)
		// return err
	}
	filterData := MyFilter{StreamTopic: customTopic, LevelList: *filterReq.Level}

	// filterData := MyFilter{StreamTopic: customTopic, LevelList: *filterReq.Level, SourceList: *filterReq.Platform}
	log.DebugLogger.Println(filterData)

	json, err := json.Marshal(filterData)
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}

	err = client.Set(getUserFilterKey(filterReq.ProfileId), json, 0).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	res := client.Get(getUserFilterKey(filterReq.ProfileId))
	log.InfoLogger.Println(res.Args()...)

	return nil
}

func GetFilterByUser(profileId string) (string, error) {

	val, err := client.Get(getUserFilterKey(profileId)).Result()
	if err != nil {
		log.ErrorLogger.Println(err)
		return "", err
	}
	log.DebugLogger.Println(val)
	return val, nil
}

// ==========={ActiveUser : [ProfileId]}=============

const ActiveUserKey = "active_user_key"

func PushUserToActiveSet(profileId string) error {

	err := client.SAdd(ActiveUserKey, profileId).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func PopUserFromActiveSet(profileId string) error {

	err := client.SRem(ActiveUserKey, profileId).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func GetActiveUsers(client *redis.Client) string {

	res, err := client.Get(ActiveUserKey).Result()
	if err != nil {
		log.ErrorLogger.Println(err)
		return ""
	}
	return res
}

type ActiveUser struct {
	JobId string
}

// --------{ActiveUser : JobId}--------------

func getUserJobKey(profileId string) string {
	return "user_job_" + profileId
}

func AddUserJobIdPair(profileId string, jobId string) error {

	err := client.Set(getUserJobKey(profileId), jobId, 0).Err()
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func GetJobIdByUser(profileId string) string {

	res, err := client.Get(getUserJobKey(profileId)).Result()
	if err != nil {
		log.ErrorLogger.Println(err)
		return ""
	}

	return res

}
