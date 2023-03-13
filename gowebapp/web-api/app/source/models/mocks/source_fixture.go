package mocks

import (
	"fmt"
	"time"
	"web-api/app/source/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fixture struct {
	Sources []schema.Source
}

func NewFixture() *Fixture {
	return &Fixture{}
}

func (f *Fixture) SetupSources(sources []schema.Source) {
	f.Sources = sources
}

func (f *Fixture) Teardown() {
	// Clean up the test data here
}

func (f *Fixture) GetSource(id uuid.UUID) (schema.Source, error) {
	for _, source := range f.Sources {
		if source.Base.ID == id {
			return source, nil
		}
	}
	return schema.Source{}, fmt.Errorf("source not found")
}

func (f *Fixture) UpdateSource(id uuid.UUID, name string) (schema.Source, error) {
	for _, source := range f.Sources {
		if source.Base.ID == id {
			source.Name = name
			source.Base.UpdatedAt = time.Now()
		}
	}
	return schema.Source{}, fmt.Errorf("source not found")
}

func (f *Fixture) GetSourcesByProfileId(id uuid.UUID) (*[]schema.Source, error) {
	res := make([]schema.Source, 0, 10)
	for _, source := range f.Sources {
		if source.ProfileId == id {
			res = append(res, source)
		}
	}
	return &res, nil
}

func (f *Fixture) DeleteSource(id uuid.UUID) (schema.Source, error) {
	for _, source := range f.Sources {
		if source.Base.ID != id {
			source.Base.DeletedAt = gorm.DeletedAt{
				Time: time.Now(),
			}
			return source, nil
		}
	}
	return schema.Source{}, gorm.ErrRecordNotFound
}

// source1 := schema.Source{
// 	ProfileId:  uuid.New(),
// 	TeamId:     uuid.New(),
// 	Name:       "First",
// 	SourceType: 1,
// 	Base: schema.Base{
// 		ID:        uuid.New(),
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	},
// }
