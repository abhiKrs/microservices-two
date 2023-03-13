package pgsql

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"auth/db"
	schema "auth/db/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailAccessOperator struct {
	dbUrl   string
	dbModel *schema.Email
}

// type Email interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.Email, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.Email, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.Email, error)
// 	Update(ctx context.Context, Email *UpdateRequest) (*schema.Email, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewEmailAccess() *EmailAccessOperator {
	return &EmailAccessOperator{
		dbUrl:   db.DBConnectUrl(),
		dbModel: &schema.Email{},
	}
}

// db.Scopes(UserTable(user)).Create(&user)
func (eao *EmailAccessOperator) CreateEmail(email string, userId uuid.UUID, isPrimary bool) (*schema.Email, error) {
	// dbEmail := schema.Email{Email: email, Verified: false, UserId: userId, Primary: isPrimary}
	var dbEmail = schema.Email{}
	// result := eao.db.Create(&dbEmail)

	sqlCmd := fmt.Sprintf("insert into emails (user_id, email, verified, primary) values ('{%s}', '%s', %t, %t);", userId, email, false, isPrimary)

	payload := `{"operation": "exec", "metadata": {"sql": "` + sqlCmd + `" }}`

	client := http.Client{}
	// Insert order using Dapr output binding via HTTP Post
	req, err := http.NewRequest("POST", eao.dbUrl, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(res)
	dbEmail = schema.Email{}
	return &dbEmail, nil

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

	log.Println(&dbEmail)

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
		log.Println(result.Error)
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

	log.Println(&dbEmail)
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
