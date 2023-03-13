package mocks

import (
	// "context"
	"log"
	"web-api/app/source/dao"
	"web-api/app/source/models"

	// "web-api/app/source/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSourceService struct {
	Dao dao.SourceDataAccess
	mock.Mock
	GetSourcesFn    func(profileId uuid.UUID) (*models.AllSourceResponse, error)
	CreateSourceFn  func(profileId uuid.UUID, pModel models.CreateSourceRequest) (*models.SingleSourceResponse, error)
	GetSourceByIdFn func(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error)
	UpdateSourceFn  func(sourceId uuid.UUID, profileId uuid.UUID, req *models.UpdateSourceRequest) (*models.SingleSourceResponse, error)
	DeleteFn        func(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error)
}

func NewMockSourceService(dao dao.SourceDataAccess) *MockSourceService {
	return &MockSourceService{
		Dao: dao,
	}
}

func (m *MockSourceService) GetSourcesByProfileId(profileId uuid.UUID) (*models.AllSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.GetSourcesFn != nil {
		log.Print("return custom mock")
		return m.GetSourcesFn(profileId)
	}
	return &models.AllSourceResponse{IsSuccessful: true, Data: []models.SourceBody{}}, nil
	// return &models.AllSourceResponse{IsSuccessful: true, Data: []schema.Source{}}, nil
}

func (m *MockSourceService) CreateSource(profileId uuid.UUID, pModel models.CreateSourceRequest) (*models.SingleSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.CreateSourceFn != nil {
		log.Print("return custom mock")
		return m.CreateSourceFn(profileId, pModel)
	}
	return &models.SingleSourceResponse{IsSuccessful: true, Data: models.SourceBody{}}, nil
	// return &models.SingleSourceResponse{IsSuccessful: true, Data: schema.Source{}}, nil
}

func (m *MockSourceService) GetSourceById(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.GetSourceByIdFn != nil {
		log.Print("return custom mock")
		return m.GetSourceByIdFn(sourceId, profileId)
	}
	return &models.SingleSourceResponse{IsSuccessful: true, Data: models.SourceBody{}}, nil
	// return &models.SingleSourceResponse{IsSuccessful: true, Data: schema.Source{}}, nil
}

func (m *MockSourceService) Delete(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.DeleteFn != nil {
		log.Print("return custom mock")
		return m.DeleteFn(sourceId, profileId)
	}
	return &models.SingleSourceResponse{IsSuccessful: true, Data: models.SourceBody{}}, nil
	// return &models.SingleSourceResponse{IsSuccessful: true, Data: schema.Source{}}, nil
}

func (m *MockSourceService) UpdateSource(sourceId uuid.UUID, profileId uuid.UUID, req *models.UpdateSourceRequest) (*models.SingleSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.UpdateSourceFn != nil {
		log.Print("return custom mock")
		return m.UpdateSourceFn(profileId, sourceId, req)
	}
	return &models.SingleSourceResponse{IsSuccessful: true, Data: models.SourceBody{}}, nil
	// return &models.SingleSourceResponse{IsSuccessful: true, Data: schema.Source{}}, nil
}
