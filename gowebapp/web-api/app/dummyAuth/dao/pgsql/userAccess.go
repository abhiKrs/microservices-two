package pgsql

import (
	"net/url"

	"web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/schema"
	"web-api/app/utility/filter"

	"gorm.io/gorm"
)

// ------------------------------------------------------------------
type Filter struct {
	Base filter.Filter

	Role      string `json:"role"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func Filters(queries url.Values) *Filter {
	f := filter.New(queries)
	if queries.Has("firstName") || queries.Has("lastName") {
		f.Search = true
	}
	return &Filter{
		Base: *f,

		Role:      queries.Get("role"),
		FirstName: queries.Get("firstName"),
		LastName:  queries.Get("lastName"),
	}
}

// ---------------------------------------------------------------------

type UserAccessOperator struct {
	db *gorm.DB
}

// type User interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.User, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.User, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.User, error)
// 	Update(ctx context.Context, user *UpdateRequest) (*schema.User, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewUserAccess(db *gorm.DB) *UserAccessOperator {
	return &UserAccessOperator{
		db: db,
	}
}

func (uao *UserAccessOperator) Create(userRole constants.DBUserRole) (*schema.User, error) {
	dbUser := schema.User{Role: uint8(constants.Users)}

	result := uao.db.Create(&dbUser)

	return &dbUser, result.Error

}
