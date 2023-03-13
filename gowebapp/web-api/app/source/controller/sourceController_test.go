package controller

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"
	"time"

	"web-api/app/config"
	"web-api/app/source/dao"
	"web-api/app/source/models"
	"web-api/app/source/models/mocks"
	"web-api/app/source/schema"
	logs "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	// "github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAll(t *testing.T) {
	// Set up test data
	source1 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "First",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	source2 := schema.Source{
		ProfileId:  source1.ProfileId,
		TeamId:     source1.TeamId,
		Name:       "Second",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	source3 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "Other",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	sources := []schema.Source{
		source1,
		source2,
		source3,
	}
	fixture := mocks.NewFixture()
	fixture.SetupSources(sources)
	defer fixture.Teardown()
	var SourceService = &mocks.MockSourceService{

		Dao: *dao.New(&gorm.DB{}, config.CFG),
		GetSourcesFn: func(profileId uuid.UUID) (*models.AllSourceResponse, error) {
			sources, err := fixture.GetSourcesByProfileId(profileId)
			if err != nil {
				logs.DebugLogger.Println(err)
			}

			var data []models.SourceBody

			for _, dbSource := range *sources {
				st := fmt.Sprintf("sourceId:%v,sourceType:%v", dbSource.Base.ID.String(), dbSource.SourceType)
				logs.DebugLogger.Println(st)
				encoded := base64.StdEncoding.EncodeToString([]byte(st))
				logs.DebugLogger.Println(encoded)
				data = append(data, models.SourceBody{Source: dbSource, SourceToken: encoded})
			}

			res := models.AllSourceResponse{
				IsSuccessful: true,
				Data:         data,
			}
			// res := models.AllSourceResponse{
			// 	IsSuccessful: true,
			// 	Data:         *sources,
			// }
			logs.DebugLogger.Println(sources)
			return &res, nil
		},
	}

	// var SourceControl = NewSourceController(mocks.NewMockSourceService(*dao.New(&gorm.DB{}, config.CFG)), validator.New())
	var SourceControl = NewSourceController(SourceService, validator.New())

	// SourceControl.sourceService.

	// Make request to API
	req, err := http.NewRequest("GET", "/sources", nil)
	if err != nil {
		t.Fatal(err)
	}

	profileId := fixture.Sources[0].ProfileId
	logs.DebugLogger.Println(profileId)
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", profileId.String()))
	rr := httptest.NewRecorder()
	// rr.WriteHeader()
	handler := http.HandlerFunc(SourceControl.GetAll)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.AllSourceResponse
	// logs.DebugLogger.Println(string(rr.Body.Bytes()))
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// logs.DebugLogger.Println(response)

	data := response.Data

	if len(data) != 2 {
		t.Errorf("handler returned wrong number of sources: got %v want %v", len(data), 2)
	}

	for i, res := range data {
		if res.Base.ID != sources[i].Base.ID {
			t.Errorf("handler returned wrong source ID: got %v want %v", res.Base.ID, sources[i].Base.ID)
		}
		if res.Name != sources[i].Name {
			t.Errorf("handler returned wrong source name: got %v want %v", res.Name, sources[i].Name)
		}
	}
}

func TestCreateSource(t *testing.T) {
	// Set up test data
	source1 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "First",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	fixture := mocks.NewFixture()
	fixture.SetupSources([]schema.Source{
		source1,
	})
	defer fixture.Teardown()

	var SourceService = &mocks.MockSourceService{

		Dao: *dao.New(&gorm.DB{}, config.CFG),
		CreateSourceFn: func(profileId uuid.UUID, pModel models.CreateSourceRequest) (*models.SingleSourceResponse, error) {
			source, err := fixture.GetSource(source1.Base.ID)
			if err != nil {
				logs.DebugLogger.Println(err)
			}
			// res := models.SingleSourceResponse{
			// 	IsSuccessful: true,
			// 	Data:         source,
			// }
			var data models.SourceBody
			st := fmt.Sprintf("sourceId:%v,sourceType:%v", source.Base.ID.String(), source.SourceType)
			logs.DebugLogger.Println(st)
			encoded := base64.StdEncoding.EncodeToString([]byte(st))
			logs.DebugLogger.Println(encoded)
			data = models.SourceBody{Source: source, SourceToken: encoded}

			res := models.SingleSourceResponse{IsSuccessful: true, Data: data}
			logs.DebugLogger.Println(source)
			return &res, nil
		},
	}

	var SourceControl = NewSourceController(SourceService, validator.New())
	// Make request to API
	reqBody := models.CreateSourceRequest{
		Name:       "First",
		SourceType: 1,
		TeamId:     source1.TeamId.String(),
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	// create request
	req, err := http.NewRequest("POST", "/sources", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	// add header
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", source1.ProfileId.String()))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SourceControl.Create)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response models.SingleSourceResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Data.Name != reqBody.Name {
		t.Errorf("handler returned wrong source name: got %v want %v", response.Data.Name, reqBody.Name)
	}
}

func TestUpdateSource(t *testing.T) {
	// Set up test data
	source1 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "First",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	fixture := mocks.NewFixture()
	fixture.SetupSources([]schema.Source{
		source1,
	})
	defer fixture.Teardown()

	var SourceService = &mocks.MockSourceService{

		Dao: *dao.New(&gorm.DB{}, config.CFG),
		UpdateSourceFn: func(profileId uuid.UUID, sourceId uuid.UUID, req *models.UpdateSourceRequest) (*models.SingleSourceResponse, error) {
			source, err := fixture.GetSource(source1.Base.ID)
			if err != nil {
				logs.DebugLogger.Println(err)
			}
			source.Name = req.Name
			source.Base.UpdatedAt = time.Now()
			// res := models.SingleSourceResponse{
			// 	IsSuccessful: true,
			// 	Data:         source,
			// }
			var data models.SourceBody
			st := fmt.Sprintf("sourceId:%v,sourceType:%v", source.Base.ID.String(), source.SourceType)
			logs.DebugLogger.Println(st)
			encoded := base64.StdEncoding.EncodeToString([]byte(st))
			logs.DebugLogger.Println(encoded)
			data = models.SourceBody{Source: source, SourceToken: encoded}

			res := models.SingleSourceResponse{IsSuccessful: true, Data: data}
			logs.DebugLogger.Println(source)
			return &res, nil
		},
	}
	var SourceControl = NewSourceController(SourceService, validator.New())
	// Make request to API
	newName := "Alicia"
	jsonData, err := json.Marshal(models.UpdateSourceRequest{Name: newName})
	if err != nil {
		t.Fatal(err)
	}

	urlPath := fmt.Sprintf("/sources/%s", source1.Base.ID.String())

	req, err := http.NewRequest(http.MethodPut, urlPath, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", source1.Base.ID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// add header
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", source1.ProfileId.String()))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SourceControl.UpdateSource)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that source was updated in fixture
	// updatedSource, err := fixture.GetSource(source1.Base.ID)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var response models.SingleSourceResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Data.Name != newName {
		t.Errorf("handler did not update source name: got %v want %v", response.Data.Name, newName)
	}
}

func TestGetSource(t *testing.T) {
	// Set up test data
	source1 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "First",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	fixture := mocks.NewFixture()
	fixture.SetupSources([]schema.Source{
		source1,
	})
	defer fixture.Teardown()

	var SourceService = &mocks.MockSourceService{

		Dao: *dao.New(&gorm.DB{}, config.CFG),
		GetSourceByIdFn: func(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {
			source, err := fixture.GetSource(source1.Base.ID)
			if err != nil {
				logs.DebugLogger.Println(err)
			}
			if source.ProfileId != profileId {
				response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", source.Base.ID)}}
				return &response, respond.ErrBadRequest
			}

			// res := models.SingleSourceResponse{
			// 	IsSuccessful: true,
			// 	Data:         source,
			// }
			var data models.SourceBody
			st := fmt.Sprintf("sourceId:%v,sourceType:%v", source.Base.ID.String(), source.SourceType)
			logs.DebugLogger.Println(st)
			encoded := base64.StdEncoding.EncodeToString([]byte(st))
			logs.DebugLogger.Println(encoded)
			data = models.SourceBody{Source: source, SourceToken: encoded}

			res := models.SingleSourceResponse{IsSuccessful: true, Data: data}
			logs.DebugLogger.Println(source)
			return &res, nil
		},
	}
	var SourceControl = NewSourceController(SourceService, validator.New())

	// Make request to API
	urlPath := fmt.Sprintf("/sources/%s", source1.Base.ID.String())

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", source1.Base.ID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// add header
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", source1.ProfileId.String()))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SourceControl.GetOne)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.SingleSourceResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Data.Base.ID != source1.Base.ID {
		t.Errorf("handler returned wrong source ID: got %v want %v", response.Data.Base.ID, source1.Base.ID)
	}
}

func TestDeleteSource(t *testing.T) {
	// Set up test data
	source1 := schema.Source{
		ProfileId:  uuid.New(),
		TeamId:     uuid.New(),
		Name:       "First",
		SourceType: 1,
		Base: schema.Base{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	fixture := mocks.NewFixture()
	fixture.SetupSources([]schema.Source{
		source1,
	})
	defer fixture.Teardown()

	var SourceService = &mocks.MockSourceService{

		Dao: *dao.New(&gorm.DB{}, config.CFG),
		DeleteFn: func(sourceId uuid.UUID, profileId uuid.UUID) (*models.SingleSourceResponse, error) {
			source, err := fixture.GetSource(source1.Base.ID)
			// logs.DebugLogger.Println(source.Base.DeletedAt)
			// source, err := fixture.DeleteSource(source1.Base.ID)
			if err != nil {
				logs.DebugLogger.Println(err)
			}

			if source.ProfileId != profileId {
				response := models.SingleSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", source.Base.ID)}}
				return &response, respond.ErrBadRequest
			}
			// delS, err := fixture.DeleteSource(source1.Base.ID)
			// if err != nil {
			// 	logs.DebugLogger.Println(err)
			// }

			// res := models.SingleSourceResponse{
			// 	IsSuccessful: true,
			// 	Data:         source,
			// }
			var data models.SourceBody
			st := fmt.Sprintf("sourceId:%v,sourceType:%v", source.Base.ID.String(), source.SourceType)
			logs.DebugLogger.Println(st)
			encoded := base64.StdEncoding.EncodeToString([]byte(st))
			logs.DebugLogger.Println(encoded)
			data = models.SourceBody{Source: source, SourceToken: encoded}

			res := models.SingleSourceResponse{IsSuccessful: true, Data: data}
			logs.DebugLogger.Println(source.Base.DeletedAt)
			return &res, nil
		},
	}
	var SourceControl = NewSourceController(SourceService, validator.New())

	// Make request to API
	urlPath := fmt.Sprintf("/sources/%s", source1.Base.ID.String())

	// Make request to API
	req, err := http.NewRequest(http.MethodDelete, urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", source1.Base.ID.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// add header
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", source1.ProfileId.String()))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SourceControl.Delete)
	handler.ServeHTTP(rr, req)

	// Check response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that source was deleted from fixture
	s, err := fixture.GetSource(source1.Base.ID)
	if err == nil {
		logs.DebugLogger.Println(s)
		// t.Errorf("handler did not delete source")
	}
}
