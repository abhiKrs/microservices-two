package pgsql

import (
	"log"

	// "web-api/app/config"
	schema "auth/db/schema"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessRequestAccessOperator struct {
	// db      *gorm.DB
	dbModel *schema.AccessRequest
}

// type AccessRequest interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.AccessRequest, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.AccessRequest, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.AccessRequest, error)
// 	Update(ctx context.Context, AccessRequest *UpdateRequest) (*schema.AccessRequest, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewAccessRequestAccess() *AccessRequestAccessOperator {
	return &AccessRequestAccessOperator{
		// db: db,
	}
}

// db.Scopes(profileTable(profile)).Create(&profile)
func (eao *AccessRequestAccessOperator) CreateAccessRequest(profileId uuid.UUID, emailId uuid.UUID) (*schema.AccessRequest, error) {
	dbAccessRequest := schema.AccessRequest{ProfileId: profileId, Approved: false, EmailId: emailId}
	// var db *gorm.DB
	result := eao.db.Create(&dbAccessRequest)

	if result.Error != nil {
		return &dbAccessRequest, result.Error
	}
	return &dbAccessRequest, nil

}

func (eao *AccessRequestAccessOperator) GetAccessRequestbyId(id uuid.UUID) (*schema.AccessRequest, error) {
	dbAccessRequest := schema.AccessRequest{}

	result := eao.db.First(dbAccessRequest, id)

	return &dbAccessRequest, result.Error
}

func (eao *AccessRequestAccessOperator) FindAccessRequestByEmailId(emailId uuid.UUID) (*schema.AccessRequest, error) {
	// dbAccessRequest := schema.AccessRequest{ProfileId: profileId, RequestType: requestType, AccessLevel: accessLevel}
	var dbAccessRequests schema.AccessRequest

	result := eao.db.Where(&schema.AccessRequest{EmailId: emailId}).Find(&dbAccessRequests)

	log.Println(&dbAccessRequests)

	// return &dbAccessRequest, result.Error
	if result.Error != nil {
		log.Println(result.Error)
		return &dbAccessRequests, result.Error
	}
	if result.RowsAffected < 1 {
		log.Println("empty results")
		return &dbAccessRequests, gorm.ErrRecordNotFound
	}

	return &dbAccessRequests, nil
}

func (eao *AccessRequestAccessOperator) FindAccessRequestByProfileId(profileId uuid.UUID) (*schema.AccessRequest, error) {
	// dbAccessRequest := schema.AccessRequest{ProfileId: profileId, RequestType: requestType, AccessLevel: accessLevel}
	var dbAccessRequests schema.AccessRequest

	result := eao.db.Where(&schema.AccessRequest{ProfileId: profileId}).First(&dbAccessRequests)

	log.Println(&dbAccessRequests)

	// return &dbAccessRequest, result.Error
	if result.Error != nil {
		log.Println(result.Error)
		return &dbAccessRequests, result.Error
	}
	if result.RowsAffected < 1 {
		log.Println("empty results")
		return &dbAccessRequests, gorm.ErrRecordNotFound
	}

	return &dbAccessRequests, nil
}

func (eao *AccessRequestAccessOperator) DeleteAccessRequest(id uuid.UUID) (*schema.AccessRequest, error) {
	dbAccessRequest := schema.AccessRequest{}

	result := eao.db.Delete(dbAccessRequest, id)

	if result.Error != nil {
		return &dbAccessRequest, result.Error
	}
	return &dbAccessRequest, nil

}

func (eao *AccessRequestAccessOperator) ApproveAccessRequest(dbModel *schema.AccessRequest) (*schema.AccessRequest, error) {

	result := eao.db.Model(dbModel).Updates(&schema.AccessRequest{Approved: true})

	if result.Error != nil {
		return dbModel, result.Error
	}
	return dbModel, nil

}

func (eao *AccessRequestAccessOperator) RejectAccessRequest(dbModel *schema.AccessRequest) (*schema.AccessRequest, error) {

	result := eao.db.Delete(dbModel)

	if result.Error != nil {
		return dbModel, result.Error
	}
	return dbModel, nil

}
