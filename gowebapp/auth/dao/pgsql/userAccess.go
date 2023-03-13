package pgsql

import (
	// "context"
	// "encoding/json"
	// "io"
	// "log"
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
	db      *gorm.DB
	dbModel *schema.User
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

func (uao *UserAccessOperator) Create(userRole constants.UserRole) (*schema.User, error) {
	var dbUser schema.User
	dbUser = schema.User{Role: uint8(constants.Users)}

	result := uao.db.Create(&dbUser)

	return &dbUser, result.Error

}

// func (u *AuthorUseCase) Create(ctx context.Context, r *author.CreateRequest) (*schema.User, error) {
// 	return u.repo.Create(ctx, r)
// }

// func (u *UserAccessOperator) Create(pModel *CreateRequest) (*schema.User, error) {
// 	var dbUser schema.User

// 	*u.dbModel = schema.User{*pModel}

// 	// var db *gorm.DB
// 	result := u.db.Create(&dbUser)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return dbUser, nil

// }

// func (u *UserAccessOperator) GetUser(Id uuid.UUID) (*schema.User, error) {
// 	dbUser := schema.User{}
// 	var dbEmail schema.Email
// 	u.db.First(&dbEmail, Id)

// 	result := u.db.First(&dbUser, dbEmail.UserId)
// 	log.Println(&dbUser)
// 	if result.Error != nil {
// 		return dbUser, result.Error
// 	}
// 	return dbUser, nil
// }

// func (u *UserAccessOperator) FindUserByEmailId(emailId uuid.UUID) (*schema.User, error) {
// 	dbUser := schema.User{}
// 	var dbEmail schema.Email
// 	u.db.First(&dbEmail, emailId)

// 	result := u.db.First(&dbUser, dbEmail.UserId)
// 	log.Println(&dbUser)
// 	if result.Error != nil {
// 		return dbUser, result.Error
// 	}
// 	return dbUser, nil
// }

// func FindCredentialsBypUIdandAuthType(pUid string, authType uint8, db *gorm.DB) (schema.Credentials, error) {
// 	// var dbCredentials schema.Credentials

// 	dbCredentials := schema.Credentials{}
// 	result := db.Where(&schema.Credentials{PUid: pUid, AuthType: authType}).First(&dbCredentials)

// 	if result.RowsAffected > 1 {
// 		log.Println(" found multiple entries for unique key")
// 		return
// 	}
// 	log.Println(&dbCredentials)
// 	if result.Error != nil {
// 		return dbCredentials, result.Error
// 	}
// 	return dbCredentials, nil
// }

// func (u *User) BeforeSave() error {
// 	hashedPassword, err := Hash(u.Password)
// 	if err != nil {
// 		return err
// 	}
// 	u.Password = string(hashedPassword)
// 	return nil
// }

// func (u *User) Prepare() {
// 	u.ID = 0
// 	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
// 	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
// 	u.CreatedAt = time.Now()
// 	u.UpdatedAt = time.Now()
// }

// func (u *User) Validate(action string) error {
// 	switch strings.ToLower(action) {
// 	case "update":
// 		if u.Nickname == "" {
// 			return errors.New("Required Nickname")
// 		}
// 		if u.Password == "" {
// 			return errors.New("Required Password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("Required Email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("Invalid Email")
// 		}

// 		return nil
// 	case "login":
// 		if u.Password == "" {
// 			return errors.New("Required Password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("Required Email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("Invalid Email")
// 		}
// 		return nil

// 	default:
// 		if u.Nickname == "" {
// 			return errors.New("Required Nickname")
// 		}
// 		if u.Password == "" {
// 			return errors.New("Required Password")
// 		}
// 		if u.Email == "" {
// 			return errors.New("Required Email")
// 		}
// 		if err := checkmail.ValidateFormat(u.Email); err != nil {
// 			return errors.New("Invalid Email")
// 		}
// 		return nil
// 	}
// }

// func (u *User) SaveUser(db *gorm.DB) (*User, error) {

// 	var err error
// 	err = db.Debug().Create(&u).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	return u, nil
// }

// func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
// 	var err error
// 	users := []User{}
// 	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
// 	if err != nil {
// 		return &[]User{}, err
// 	}
// 	return &users, err
// }

// func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
// 	var err error
// 	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	if gorm.IsRecordNotFoundError(err) {
// 		return &User{}, errors.New("User Not Found")
// 	}
// 	return u, err
// }

// func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

// 	// To hash the password
// 	err := u.BeforeSave()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
// 		map[string]interface{}{
// 			"password":   u.Password,
// 			"nickname":   u.Nickname,
// 			"email":      u.Email,
// 			"updated_at": time.Now(),
// 		},
// 	)
// 	if db.Error != nil {
// 		return &User{}, db.Error
// 	}
// 	// This is the display the updated user
// 	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
// 	if err != nil {
// 		return &User{}, err
// 	}
// 	return u, nil
// }

// func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

// 	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

// 	if db.Error != nil {
// 		return 0, db.Error
// 	}
// 	return db.RowsAffected, nil
// }
