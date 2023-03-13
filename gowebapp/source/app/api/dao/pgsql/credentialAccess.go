package pgsql

import (
	constants "web-api/app/dummyAuth/constants"
	"web-api/app/dummyAuth/schema"
	log "web-api/app/utility/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CredentialAccessOperator struct {
	db *gorm.DB
}

// type Credential interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.Credential, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.Credential, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.Credential, error)
// 	Update(ctx context.Context, Credential *UpdateRequest) (*schema.Credential, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewcredentialAccess(db *gorm.DB) *CredentialAccessOperator {
	return &CredentialAccessOperator{
		db: db,
	}
}

func (cao *CredentialAccessOperator) CreateGoogleCredential(
	userId uuid.UUID,
	identifier string,
	authType constants.DBAuthType,
	providerType constants.DBProviderType,
) (*schema.Credential, error) {

	log.InfoLogger.Println("creating credential")
	dbCredential := schema.Credential{
		UserId:       userId,
		AuthType:     uint8(authType),
		ProviderType: uint8(providerType),
		Identifier:   identifier,
	}
	// var db *gorm.DB
	result := cao.db.Create(&dbCredential)

	if result.Error != nil {
		return &dbCredential, result.Error
	}
	return &dbCredential, nil
}

func (cao *CredentialAccessOperator) CreatePasswordCredential(
	userId uuid.UUID,
	hashPass *string,
	hashAlgo constants.DBPasswordAlgo,
	identifier string,
	authType constants.DBAuthType,
	providerType constants.DBProviderType,
) (*schema.Credential, error) {

	log.InfoLogger.Println("creating credential")
	algoPass := hashAlgo.String()
	dbCredential := schema.Credential{
		UserId:       userId,
		AuthType:     uint8(authType),
		ProviderType: uint8(providerType),
		PasswordHash: hashPass,
		Identifier:   identifier,
		PasswordAlgo: &algoPass,
	}
	// var db *gorm.DB
	result := cao.db.Create(&dbCredential)

	if result.Error != nil {
		return &dbCredential, result.Error
	}
	return &dbCredential, nil
}

func (cao *CredentialAccessOperator) GetCredential(userId uuid.UUID, aType constants.DBAuthType, pType constants.DBProviderType) (*schema.Credential, error) {

	var dbCredential schema.Credential

	result := cao.db.Where(schema.Credential{UserId: userId, AuthType: uint8(aType), ProviderType: uint8(pType)}).First(&dbCredential)

	if result.Error != nil {
		log.ErrorLogger.Println(result.Error)
		return &dbCredential, result.Error
	}
	if result.RowsAffected < 1 {
		log.DebugLogger.Println("empty results")
		return &dbCredential, gorm.ErrRecordNotFound
	}

	return &dbCredential, nil

}

func (cao *CredentialAccessOperator) GetCredentialbyIdentifier(identifier string, aType constants.DBAuthType, pType constants.DBProviderType) (*schema.Credential, error) {

	var dbCredential schema.Credential

	result := cao.db.Where(schema.Credential{Identifier: identifier, AuthType: uint8(aType), ProviderType: uint8(pType)}).First(&dbCredential)

	if result.Error != nil {
		log.ErrorLogger.Println(result.Error)
		return &dbCredential, result.Error
	}
	if result.RowsAffected < 1 {
		log.DebugLogger.Println("empty results")
		return &dbCredential, gorm.ErrRecordNotFound
	}

	return &dbCredential, nil

}

func (cao *CredentialAccessOperator) UpdatePasswordCredential(dbModel *schema.Credential, hashPassword string, passwordAlgo constants.DBPasswordAlgo) (*schema.Credential, error) {

	passAlgo := passwordAlgo.String()

	result := cao.db.Model(dbModel).Updates(&schema.Credential{PasswordHash: &hashPassword, PasswordAlgo: &passAlgo})

	if result.Error != nil {
		log.ErrorLogger.Println(result.Error)
		return dbModel, result.Error
	}
	if result.RowsAffected < 1 {
		log.DebugLogger.Println("zero updates")
		return dbModel, gorm.ErrRecordNotFound
	}

	return dbModel, nil

}
