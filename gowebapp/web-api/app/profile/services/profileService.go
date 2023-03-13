package services

import (
	"web-api/app/dummyAuth/models"
	"web-api/app/dummyAuth/schema"
	"web-api/app/profile/dao"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
)

type ProfileService struct {
	dao dao.ProfileDataAccess
}

func NewProfileService(dao dao.ProfileDataAccess) *ProfileService {
	return &ProfileService{
		dao: dao,
	}
}

func (ps *ProfileService) GetUserProfileById(profileId uuid.UUID) (*models.ProfileResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get Profile
	dbProfile, err := ps.dao.ProfileAccess.GetById(profileId)
	if err != nil {
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// Get Email(primary)
	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbProfile.UserId)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// Get AccessRequest
	dbAccessRequest, err := ps.dao.AccessRequestAccess.FindAccessRequestByEmailId(dbEmail.ID)
	if err != nil {
		log.ErrorLogger.Println(err)

	}

	bodyPayload := models.ProfileBodyResponse{
		Profile:        *dbProfile,
		Email:          &dbEmail.Email,
		AccessApproved: dbAccessRequest.Approved,
	}

	response := models.ProfileResponse{IsSuccessful: true, UserBody: bodyPayload}
	return &response, nil

}

func (ps *ProfileService) UpdateProfile(profileId uuid.UUID, req *schema.ProfileBase) (*models.ProfileResponse, error) {

	var err error

	// Get Profile
	dbProfile, err := ps.dao.ProfileAccess.GetById(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	// Get Email(Primary)
	dbEmail, err := ps.dao.EmailAccess.FindPrimaryEmailbyUserId(dbProfile.UserId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}
	// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
	if !dbProfile.Onboarded {
		// Create Access Request
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"user must be onboarded first, to update profile"}}
		return &response, respond.ErrInvalidRequest
	}
	// -----------------------------------------------------------------------------------------------------
	// Get Access Request
	dbAccessRequest, err := ps.dao.AccessRequestAccess.FindAccessRequestByProfileId(dbProfile.ID)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
		return &response, respond.ErrInvalidRequest
	}

	// Update Profile
	if dbAccessRequest.Approved {
		dbProfile, err = ps.dao.ProfileAccess.Update(dbProfile, req)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.ProfileResponse{IsSuccessful: false, Message: []string{err.Error()}}
			return &response, err
		}
	} else {
		log.DebugLogger.Println("Access Request not approved")
		// ----------------------don't remove. Uncomment after onboard flow moves to /onboard/{id}-------------------------
		response := models.ProfileResponse{IsSuccessful: false, Message: []string{"user has no access granted"}}
		return &response, respond.ErrInvalidRequest
		// -----------------------------------------------------------------------------------------------------

	}

	bodyPayload := models.ProfileBodyResponse{
		Profile:        *dbProfile,
		Email:          &dbEmail.Email,
		AccessApproved: dbAccessRequest.Approved,
	}
	response := models.ProfileResponse{IsSuccessful: true, UserBody: bodyPayload}
	return &response, nil

}
