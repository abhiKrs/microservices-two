package services

// import (
// 	"fmt"
// 	"web-api/app/team/models"

// 	// "web-api/app/team/schema"
// 	"web-api/app/team/dao"
// 	log "web-api/app/utility/logger"

// 	"web-api/app/utility/respond"

// 	"github.com/google/uuid"
// )

// type TeamService struct {
// 	dao dao.TeamDataAccess
// }

// func NewTeamService(dao dao.TeamDataAccess) *TeamService {
// 	return &TeamService{
// 		dao: dao,
// 	}
// }

// func (ps *TeamService) GetTeamsByProfileId(profileId uuid.UUID) (*models.GetAllTeamResponse, error) {

// 	var err error

// 	// Get Team
// 	dbTeams, err := ps.dao.TeamAccess.GetByProfileId(profileId)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		response := models.GetAllTeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// 		return &response, err
// 	}

// 	var dataList []models.TeamBodyResponse

// 	for _, dbTeam := range *dbTeams {
// 		data := models.TeamBodyResponse{
// 			TeamId:    dbTeam.Base.ID,
// 			ProfileId:   dbTeam.ProfileId,
// 			Name:        dbTeam.Name,
// 			TeamType:  dbTeam.TeamType,
// 			TeamId:      dbTeam.TeamId,
// 			TeamToken: dbTeam.Base.ID.String(),
// 		}
// 		dataList = append(dataList, data)
// 	}

// 	bodyPayload := models.GetAllTeamResponse{
// 		IsSuccessful: true,
// 		Data:         dataList,
// 	}

// 	return &bodyPayload, nil

// }

// func (ps *TeamService) CreateTeam(pModel models.CreateTeamRequest, profileId uuid.UUID) (*models.CreateTeamResponse, error) {

// 	var err error
// 	var pId uuid.UUID
// 	if pModel.TeamId != "" {
// 		pId = uuid.New()
// 	} else {
// 		pId = uuid.MustParse(pModel.TeamId)
// 	}

// 	// Get Team
// 	dbTeam, err := ps.dao.TeamAccess.Create(profileId, pModel.Name, pModel.TeamType, pId)
// 	if err != nil {
// 		response := models.CreateTeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// 		return &response, err
// 	}

// 	bodyPayload := models.TeamBodyResponse{
// 		TeamId:    dbTeam.Base.ID,
// 		ProfileId:   dbTeam.ProfileId,
// 		Name:        dbTeam.Name,
// 		TeamType:  dbTeam.TeamType,
// 		TeamId:      dbTeam.TeamId,
// 		TeamToken: dbTeam.Base.ID.String(),
// 	}

// 	response := models.CreateTeamResponse{IsSuccessful: true, Data: bodyPayload}
// 	return &response, nil
// }

// func (ps *TeamService) GetTeamById(teamId uuid.UUID, profileId uuid.UUID) (*models.GetTeamResponse, error) {

// 	// var dbEmail schema.Email
// 	var err error

// 	// Get Team
// 	dbTeam, err := ps.dao.TeamAccess.GetById(teamId)
// 	if err != nil {
// 		response := models.GetTeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// 		return &response, err
// 	}

// 	if dbTeam.ProfileId != profileId {
// 		log.ErrorLogger.Printf("Unathorised access to team: %v from different profile: %v /n", dbTeam.Base.ID, profileId)
// 		response := models.GetTeamResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to team: %v from different profile", dbTeam.Base.ID)}}
// 		return &response, respond.ErrBadRequest
// 	}

// 	bodyPayload := models.TeamBodyResponse{
// 		TeamId:    dbTeam.Base.ID,
// 		ProfileId:   dbTeam.ProfileId,
// 		Name:        dbTeam.Name,
// 		TeamType:  dbTeam.TeamType,
// 		TeamId:      dbTeam.TeamId,
// 		TeamToken: dbTeam.Base.ID.String(),
// 	}

// 	response := models.GetTeamResponse{IsSuccessful: true, Data: bodyPayload}
// 	return &response, nil

// }

// // func (ps *TeamService) UpdateTeam(profileId uuid.UUID, req *schema.TeamBase) (*models.TeamResponse, error) {

// // 	var err error

// // 	// Get Team
// // 	dbTeam, err := ps.dao.TeamAccess.GetById(profileId)
// // 	if err != nil {
// // 		log.ErrorLogger.Println(err)
// // 		response := models.TeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// // 		return &response, err
// // 	}

// // 	// Get Email(Primary)
// // 	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbTeam.UserId)
// // 	if err != nil {
// // 		log.ErrorLogger.Println(err)
// // 		response := models.TeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// // 		return &response, err
// // 	}
// // 	// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
// // 	if !dbTeam.Onboarded {
// // 		// Create Access Request
// // 		response := models.TeamResponse{IsSuccessful: false, Message: []string{"user must be onboarded first, to update profile"}}
// // 		return &response, respond.ErrInvalidRequest
// // 	}
// // 	// -----------------------------------------------------------------------------------------------------
// // 	// Get Access Request
// // 	dbAccessRequest, err := ps.dao.AccessRequestAccess.FindAccessRequestByTeamId(dbTeam.ID)
// // 	if err != nil {
// // 		log.ErrorLogger.Println(err)
// // 		response := models.TeamResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
// // 		return &response, respond.ErrInvalidRequest
// // 	}

// // 	// Update Team
// // 	if dbAccessRequest.Approved {
// // 		dbTeam, err = ps.dao.TeamAccess.Update(dbTeam, req)
// // 		if err != nil {
// // 			log.ErrorLogger.Println(err)
// // 			response := models.TeamResponse{IsSuccessful: false, Message: []string{err.Error()}}
// // 			return &response, err
// // 		}
// // 	} else {
// // 		log.DebugLogger.Println("Access Request not approved")
// // 		// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
// // 		response := models.TeamResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
// // 		return &response, respond.ErrInvalidRequest
// // 		// -----------------------------------------------------------------------------------------------------

// // 	}

// // 	bodyPayload := models.TeamBodyResponse{
// // 		Team:        *dbTeam,
// // 		Email:          &dbEmail.Email,
// // 		AccessApproved: dbAccessRequest.Approved,
// // 	}
// // 	response := models.TeamResponse{IsSuccessful: true, UserBody: bodyPayload}
// // 	return &response, nil

// // }
