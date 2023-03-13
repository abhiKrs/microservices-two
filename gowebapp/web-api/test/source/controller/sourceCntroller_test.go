package controller

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"web-api/app/source/schema"

// 	"github.com/google/uuid"
// )

// // var Source schema.Source

// // var SC = server.Server.Source()
// type Fixture struct {
// 	Sources []schema.Source
// }

// func NewFixture() *Fixture {
// 	return &Fixture{}
// }

// func (f *Fixture) SetupSources(sources []schema.Source) {
// 	f.Sources = sources
// }

// func (f *Fixture) Teardown() {
// 	// Clean up the test data here
// }

// func (f *Fixture) GetSource(id uuid.UUID) (schema.Source, error) {
// 	for _, source := range f.Sources {
// 		if source.Base.ID == id {
// 			return source, nil
// 		}
// 	}
// 	return schema.Source{}, fmt.Errorf("source not found")
// }

// func TestGetAll(t *testing.T) {
// 	// Set up test data
// 	sources := []schema.Source{
// 		{
// 			ProfileId:  uuid.New(),
// 			TeamId:     uuid.New(),
// 			Name:       "First",
// 			SourceType: 1,
// 			Base: schema.Base{
// 				ID:        uuid.New(),
// 				CreatedAt: time.Now(),
// 				UpdatedAt: time.Now(),
// 			},
// 		},
// 		{
// 			ProfileId:  uuid.New(),
// 			TeamId:     uuid.New(),
// 			Name:       "Second",
// 			SourceType: 1,
// 			Base: schema.Base{
// 				ID:        uuid.New(),
// 				CreatedAt: time.Now(),
// 				UpdatedAt: time.Now(),
// 			},
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources(sources)
// 	defer fixture.Teardown()

// 	// Make request to API
// 	req, err := http.NewRequest("GET", "/sources", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(GetAll)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	var response []schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(response) != len(sources) {
// 		t.Errorf("handler returned wrong number of sources: got %v want %v", len(response), len(sources))
// 	}

// 	for i, source := range response {
// 		if source.Base.ID != sources[i].Base.ID {
// 			t.Errorf("handler returned wrong source ID: got %v want %v", source.Base.ID, sources[i].Base.ID)
// 		}
// 		if source.Name != sources[i].Name {
// 			t.Errorf("handler returned wrong source name: got %v want %v", source.Name, sources[i].Name)
// 		}
// 	}
// }

// func TestCreateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{Name: "Alice"}
// 	fixture := NewFixture()
// 	defer fixture.Teardown()

// 	// Make request to API
// 	jsonData, err := json.Marshal(source)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/sources", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Create)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}

// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }

// func TestGetSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{ID: 1, Name: "Alice"}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	req, err := http.NewRequest("GET", "/sources")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(GetOne)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if response.Base.ID != source.Base.ID {
// 		t.Errorf("handler returned wrong source ID: got %v want %v", response.Base.ID, source.Base.ID)
// 	}
// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }

// func TestUpdateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{ID: 1, Name: "Alice"}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	newName := "Alicia"
// 	jsonData, err := json.Marshal(schema.Source{Name: newName})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/sources/1", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Update)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check that source was updated in fixture
// 	updatedSource, err := fixture.GetSource(1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if updatedSource.Name != newName {
// 		t.Errorf("handler did not update source name: got %v want %v", updatedSource.Name, newName)
// 	}
// }

// func TestDeleteSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{ID: 1, Name: "Alice"}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	req, err := http.NewRequest("DELETE", "/sources/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(Delete)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check that source was deleted from fixture
// 	_, err = fixture.GetSource(1)
// 	if err == nil {
// 		t.Errorf("handler did not delete source")
// 	}
// }
