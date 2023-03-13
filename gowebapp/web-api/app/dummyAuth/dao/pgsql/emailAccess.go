package pgsql

import (
	schema "web-api/app/dummyAuth/schema"
	log "web-api/app/utility/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailAccessOperator struct {
	db *gorm.DB
}

// type Email interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.Email, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.Email, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.Email, error)
// 	Update(ctx context.Context, Email *UpdateRequest) (*schema.Email, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewEmailAccess(db *gorm.DB) *EmailAccessOperator {
	return &EmailAccessOperator{
		db: db,
	}
}

// db.Scopes(UserTable(user)).Create(&user)
func (eao *EmailAccessOperator) CreateEmail(email string, userId uuid.UUID, isPrimary bool) (*schema.Email, error) {
	dbEmail := schema.Email{Email: email, Verified: false, UserId: userId, Primary: isPrimary}
	// var db *gorm.DB
	result := eao.db.Create(&dbEmail)

	if result.Error != nil {
		return &dbEmail, result.Error
	}
	return &dbEmail, nil

}

func (eao *EmailAccessOperator) FindEmailbyId(id uuid.UUID) (*schema.Email, error) {
	dbEmail := schema.Email{}
	// var dbEmail schema.Email
	// db.First(&dbEmail, emailId)

	result := eao.db.First(&dbEmail, id)

	log.DebugLogger.Println(&dbEmail)

	return &dbEmail, result.Error
}

func (eao *EmailAccessOperator) DeleteEmail(emailId uuid.UUID) (*schema.Email, error) {
	dbEmail := schema.Email{ID: emailId}
	// var db *gorm.DB
	result := eao.db.Delete(dbEmail)

	if result.Error != nil {
		return &dbEmail, result.Error
	}
	return &dbEmail, nil

}

func (eao *EmailAccessOperator) FindEmailbyStrAdd(emailAddress string) (*schema.Email, error) {
	dbEmail := schema.Email{}
	// var dbEmail schema.Email
	// db.First(&dbEmail, emailId)

	result := eao.db.First(&dbEmail, "email = ?", emailAddress)
	if result.Error != nil {
		log.ErrorLogger.Println(result.Error)
		return &dbEmail, result.Error
	}
	if result.RowsAffected < 1 {
		return &dbEmail, gorm.ErrRecordNotFound
	}

	return &dbEmail, nil

	// return &dbEmail, result.Error
}

func (eao *EmailAccessOperator) FindPrimaryEmailbyUserId(userId uuid.UUID) (*schema.Email, error) {
	dbEmail := schema.Email{}
	// var dbEmail schema.Email
	// db.First(&dbEmail, emailId)

	result := eao.db.Where(
		&schema.Email{
			UserId:  userId,
			Primary: true,
		}).Find(&dbEmail)

	log.DebugLogger.Println(&dbEmail)
	if result.Error != nil {
		return &dbEmail, result.Error
	}
	return &dbEmail, nil
}

func (eao *EmailAccessOperator) VerifyEmail(dbModel *schema.Email) (*schema.Email, error) {

	result := eao.db.Model(dbModel).Updates(&schema.Email{Verified: true})

	if result.Error != nil {
		return dbModel, result.Error
	}
	return dbModel, nil

}
