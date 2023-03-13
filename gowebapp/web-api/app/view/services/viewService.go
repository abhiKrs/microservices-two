package services

import (
	"fmt"
	// "web-api/app/dummyAuth/schema"
	"web-api/app/view/models"
	"web-api/app/view/schema"

	// "web-api/app/view/schema"
	log "web-api/app/utility/logger"
	"web-api/app/view/dao"

	"web-api/app/utility/respond"

	"github.com/google/uuid"
)

type ViewService struct {
	dao dao.ViewDataAccess
}

func NewViewService(dao dao.ViewDataAccess) *ViewService {
	return &ViewService{
		dao: dao,
	}
}

func (ps *ViewService) GetViewsByProfileId(profileId uuid.UUID) (*models.GetAllViewResponse, error) {

	var err error

	// Get View
	// TODO - list if teams to iterate through
	team, err := ps.dao.TeamAccess.GetByCreatorId(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}
	dbViewList, err := ps.dao.ViewAccess.GetByTeamId(team.Base.ID)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	var dataList []models.ViewBodyResponse

	for _, dbView := range *dbViewList {
		// if dbView.LevelFilter != nil {

		// }
		data := models.ViewBodyResponse{
			CreateDate: dbView.Base.CreatedAt,
			CreateViewRequest: models.CreateViewRequest{
				StreamId:      dbView.StreamId.String(),
				TeamId:        dbView.TeamId.String(),
				SourcesFilter: dbView.SourcesFilter,
				LevelFilter:   dbView.LevelFilter,
				DateFilter: &models.DateInterval{
					StartDate: dbView.StartDate,
					EndDate:   dbView.EndDate,
				},
				Name: dbView.Name,
			},
		}
		dataList = append(dataList, data)
	}

	bodyPayload := models.GetAllViewResponse{
		IsSuccessful: true,
		Data:         dataList,
	}

	return &bodyPayload, nil

}

func (ps *ViewService) GetViewByTeamId(profileId uuid.UUID) (*models.GetAllViewResponse, error) {

	var err error

	// Get View
	dbView, err := ps.dao.ViewAccess.GetByTeamId(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetAllViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	var dataList []models.ViewBodyResponse

	for _, dbView := range *dbView {
		data := models.ViewBodyResponse{
			CreateDate: dbView.Base.CreatedAt,
			CreateViewRequest: models.CreateViewRequest{
				StreamId:      dbView.StreamId.String(),
				TeamId:        dbView.TeamId.String(),
				SourcesFilter: dbView.SourcesFilter,
				LevelFilter:   dbView.LevelFilter,
				DateFilter: &models.DateInterval{
					StartDate: dbView.StartDate,
					EndDate:   dbView.EndDate,
				},
				Name: dbView.Name,
			},
		}
		dataList = append(dataList, data)
	}

	bodyPayload := models.GetAllViewResponse{
		IsSuccessful: true,
		Data:         dataList,
	}

	return &bodyPayload, nil

}

func (ps *ViewService) CreateView(pModel models.CreateViewRequest, profileId uuid.UUID) (*models.CreateViewResponse, error) {

	var err error
	var streamId *uuid.UUID
	var teamId uuid.UUID
	var dbView *schema.View
	teamId, err = uuid.Parse(pModel.TeamId)
	if err != nil {
		response := models.CreateViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if pModel.StreamId != "" {
		*streamId, err = uuid.Parse(pModel.StreamId)
		if err != nil {
			response := models.CreateViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &response, err
		}
	}

	// Get View
	// if s
	dbView, err = ps.dao.ViewAccess.Create(teamId, pModel.Name, pModel.LevelFilter, pModel.SourcesFilter, pModel.DateFilter, streamId)
	if err != nil {
		response := models.CreateViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	bodyPayload := models.ViewBodyResponse{
		CreateDate: dbView.Base.CreatedAt,
		CreateViewRequest: models.CreateViewRequest{
			StreamId:      dbView.StreamId.String(),
			TeamId:        dbView.TeamId.String(),
			SourcesFilter: dbView.SourcesFilter,
			LevelFilter:   dbView.LevelFilter,
			DateFilter: &models.DateInterval{
				StartDate: dbView.StartDate,
				EndDate:   dbView.EndDate,
			},
			Name: dbView.Name,
		},
	}

	response := models.CreateViewResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil
}

func (ps *ViewService) GetViewById(viewId uuid.UUID, teamId uuid.UUID) (*models.GetViewResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get View
	dbView, err := ps.dao.ViewAccess.GetById(viewId)
	if err != nil {
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbView.TeamId != teamId {
		log.ErrorLogger.Printf("Unathorised access to view: %v from different profile: %v /n", dbView.Base.ID, teamId)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("wrong access pattern to view: %v from different profile", dbView.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	bodyPayload := models.ViewBodyResponse{
		CreateDate: dbView.Base.CreatedAt,
		CreateViewRequest: models.CreateViewRequest{
			StreamId:      dbView.StreamId.String(),
			TeamId:        dbView.TeamId.String(),
			SourcesFilter: dbView.SourcesFilter,
			LevelFilter:   dbView.LevelFilter,
			DateFilter: &models.DateInterval{
				StartDate: dbView.StartDate,
				EndDate:   dbView.EndDate,
			},
			Name: dbView.Name,
		},
	}

	response := models.GetViewResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil

}

func (ps *ViewService) Delete(viewId uuid.UUID, teamId uuid.UUID) (*models.DeleteViewResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get View
	dbView, err := ps.dao.ViewAccess.GetById(viewId)
	if err != nil {
		response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbView.TeamId != teamId {
		log.ErrorLogger.Printf("Unathorised access to view: %v from different team: %v /n", dbView.Base.ID, teamId)
		response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to view: %v from different team", dbView.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	dbView, err = ps.dao.ViewAccess.DeleteById(viewId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.DeleteViewResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to view: %v from different team", dbView.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	response := models.DeleteViewResponse{IsSuccessful: true}
	return &response, nil

}

func (ps *ViewService) UpdateView(teamId uuid.UUID, viewId uuid.UUID, req *models.UpdateViewRequest) (*models.GetViewResponse, error) {

	var err error

	// Get View
	dbView, err := ps.dao.ViewAccess.GetById(viewId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbView.TeamId != teamId {
		log.ErrorLogger.Printf("Unathorised access to view: %v from different team: %v /n", dbView.Base.ID, teamId)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to view: %v from different team", dbView.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	// updateModel := schema.

	dbView, err = ps.dao.ViewAccess.Update(dbView, *req)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetViewResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to view: %v from different team", dbView.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	bodyPayload := models.ViewBodyResponse{
		CreateDate: dbView.Base.CreatedAt,
		CreateViewRequest: models.CreateViewRequest{
			StreamId:      dbView.StreamId.String(),
			TeamId:        dbView.TeamId.String(),
			SourcesFilter: dbView.SourcesFilter,
			LevelFilter:   dbView.LevelFilter,
			DateFilter: &models.DateInterval{
				StartDate: dbView.StartDate,
				EndDate:   dbView.EndDate,
			},
			Name: dbView.Name,
		},
	}

	response := models.GetViewResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil
}

// 	// Get Email(Primary)
// 	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbView.UserId)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		response := models.ViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
// 		return &response, err
// 	}
// 	// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
// 	if !dbView.Onboarded {
// 		// Create Access Request
// 		response := models.ViewResponse{IsSuccessful: false, Message: []string{"user must be onboarded first, to update profile"}}
// 		return &response, respond.ErrInvalidRequest
// 	}
// 	// -----------------------------------------------------------------------------------------------------
// 	// Get Access Request
// 	dbAccessRequest, err := ps.dao.AccessRequestAccess.FindAccessRequestByViewId(dbView.ID)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		response := models.ViewResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
// 		return &response, respond.ErrInvalidRequest
// 	}

// 	// Update View
// 	if dbAccessRequest.Approved {
// 		dbView, err = ps.dao.ViewAccess.Update(dbView, req)
// 		if err != nil {
// 			log.ErrorLogger.Println(err)
// 			response := models.ViewResponse{IsSuccessful: false, Message: []string{err.Error()}}
// 			return &response, err
// 		}
// 	} else {
// 		log.DebugLogger.Println("Access Request not approved")
// 		// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
// 		response := models.ViewResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
// 		return &response, respond.ErrInvalidRequest
// 		// -----------------------------------------------------------------------------------------------------

// 	}

// 	bodyPayload := models.ViewBodyResponse{
// 		View:        *dbView,
// 		Email:          &dbEmail.Email,
// 		AccessApproved: dbAccessRequest.Approved,
// 	}
// 	response := models.ViewResponse{IsSuccessful: true, UserBody: bodyPayload}
// 	return &response, nil

// }
