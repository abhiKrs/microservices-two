package services

import (
	"encoding/base64"
	"fmt"

	// "strings"
	"web-api/app/source/constants"
	"web-api/app/source/models"

	// "web-api/app/source/schema"
	"web-api/app/source/dao"
	logs "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
)

type SourceService struct {
	dao dao.SourceDataAccess
}

func NewSourceService(dao dao.SourceDataAccess) *SourceService {
	return &SourceService{
		dao: dao,
	}
}

func (ps *SourceService) GetSourcesByProfileId(profileId uuid.UUID) (*models.AllSourceResponse, error) {

	var err error

	// Get Source
	dbSources, err := ps.dao.SourceAccess.GetByProfileId(profileId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.AllSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	var data []models.SourceBody

	for _, dbSource := range *dbSources {
		st := fmt.Sprintf("sourceId:%v,sourceType:%v", dbSource.Base.ID.String(), dbSource.SourceType)
		logs.DebugLogger.Println(st)
		encoded := base64.StdEncoding.EncodeToString([]byte(st))
		logs.DebugLogger.Println(encoded)
		data = append(data, models.SourceBody{Source: dbSource, SourceToken: encoded})
	}

	bodyPayload := models.AllSourceResponse{
		IsSuccessful: true,
		Data:         data,
	}

	return &bodyPayload, nil

}

func (ps *SourceService) CreateSource(profileId uuid.UUID, pModel models.CreateSourceRequest) (*models.SingleSourceResponse, error) {

	var err error
	var pId uuid.UUID
	if pModel.TeamId != "" {
		pId = uuid.New()
	} else {
		pId = uuid.MustParse(pModel.TeamId)
	}

	// Get Source
	dbSource, err := ps.dao.SourceAccess.Create(profileId, pModel.Name, constants.SourceType(pModel.SourceType), pId)
	if err != nil {
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// bodyPayload := models.SourceBodyResponse{
	// 	SourceId:    dbSource.Base.ID,
	// 	ProfileId:   dbSource.ProfileId,
	// 	Name:        dbSource.Name,
	// 	SourceType:  dbSource.SourceType,
	// 	TeamId:      dbSource.TeamId,
	// 	SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
	// 	CreatedAt:   dbSource.Base.CreatedAt,
	// }
	var data models.SourceBody
	st := fmt.Sprintf("sourceId:%v,sourceType:%v", dbSource.Base.ID.String(), dbSource.SourceType)
	logs.DebugLogger.Println(st)
	encoded := base64.StdEncoding.EncodeToString([]byte(st))
	logs.DebugLogger.Println(encoded)
	data = models.SourceBody{Source: *dbSource, SourceToken: encoded}

	response := models.SingleSourceResponse{IsSuccessful: true, Data: data}
	// response := models.SingleSourceResponse{IsSuccessful: true, Data: *dbSource}
	return &response, nil
}

func (ps *SourceService) GetSourceById(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	// bodyPayload := models.SourceBodyResponse{
	// 	SourceId:    dbSource.Base.ID,
	// 	ProfileId:   dbSource.ProfileId,
	// 	Name:        dbSource.Name,
	// 	SourceType:  dbSource.SourceType,
	// 	TeamId:      dbSource.TeamId,
	// 	SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
	// 	CreatedAt:   dbSource.Base.CreatedAt,
	// }

	var data models.SourceBody
	st := fmt.Sprintf("sourceId:%v,sourceType:%v", dbSource.Base.ID.String(), dbSource.SourceType)
	logs.DebugLogger.Println(st)
	encoded := base64.StdEncoding.EncodeToString([]byte(st))
	logs.DebugLogger.Println(encoded)
	data = models.SourceBody{Source: *dbSource, SourceToken: encoded}

	response := models.SingleSourceResponse{IsSuccessful: true, Data: data}
	// response := models.SingleSourceResponse{IsSuccessful: true, Data: *dbSource}
	return &response, nil

}

func (ps *SourceService) Delete(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	dbSource, err = ps.dao.SourceAccess.DeleteById(sourceId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	response := models.SingleSourceResponse{IsSuccessful: true}
	return &response, nil

}

func (ps *SourceService) UpdateSource(sourceId uuid.UUID, profileId uuid.UUID, req *models.UpdateSourceRequest) (*models.SingleSourceResponse, error) {

	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	dbSource, err = ps.dao.SourceAccess.Update(dbSource, req.Name)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	// bodyPayload := models.SourceBodyResponse{
	// 	SourceId:    dbSource.Base.ID,
	// 	ProfileId:   dbSource.ProfileId,
	// 	Name:        dbSource.Name,
	// 	SourceType:  dbSource.SourceType,
	// 	TeamId:      dbSource.TeamId,
	// 	SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
	// 	CreatedAt:   dbSource.Base.CreatedAt,
	// }

	var data models.SourceBody
	st := fmt.Sprintf("sourceId:%v,sourceType:%v", dbSource.Base.ID.String(), dbSource.SourceType)
	logs.DebugLogger.Println(st)
	encoded := base64.StdEncoding.EncodeToString([]byte(st))
	logs.DebugLogger.Println(encoded)
	data = models.SourceBody{Source: *dbSource, SourceToken: encoded}

	response := models.SingleSourceResponse{IsSuccessful: true, Data: data}

	// response := models.SingleSourceResponse{IsSuccessful: true, Data: *dbSource}
	return &response, nil
}
