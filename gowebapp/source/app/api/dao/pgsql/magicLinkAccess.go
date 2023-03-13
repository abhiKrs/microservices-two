package pgsql

import (
	"time"

	constants "web-api/app/dummyAuth/constants"
	schema "web-api/app/dummyAuth/schema"
	log "web-api/app/utility/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MagicLinkAccessOperator struct {
	db *gorm.DB
}

// type MagicLink interface {
// 	Create(ctx context.Context, a *CreateRequest) (*schema.MagicLink, error)
// 	List(ctx context.Context, f *Filter) ([]*schema.MagicLink, int, error)
// 	Read(ctx context.Context, ID uuid.UUID) (*schema.MagicLink, error)
// 	Update(ctx context.Context, MagicLink *UpdateRequest) (*schema.MagicLink, error)
// 	Delete(ctx context.Context, ID uuid.UUID) error
// }

func NewMagicLinkAccess(db *gorm.DB) *MagicLinkAccessOperator {
	return &MagicLinkAccessOperator{
		db: db,
	}
}

func (ml *MagicLinkAccessOperator) CreateMagicLink(linkType constants.DBMagicLinkType, emailId, userId uuid.UUID) (*schema.MagicLink, error) {
	token := uuid.New()

	dbMagicLink := schema.MagicLink{
		Token:      token,
		UserId:     userId,
		EmailId:    emailId,
		LinkType:   uint8(linkType),
		CreatedAt:  time.Now(),
		ExpiryTime: time.Now().Add(1 * time.Hour),
	}

	result := ml.db.Create(&dbMagicLink)

	if result.Error != nil {
		return &dbMagicLink, result.Error
	}
	return &dbMagicLink, nil

}

func (ml *MagicLinkAccessOperator) FindMagicLink(token uuid.UUID) (*schema.MagicLink, error) {
	var dbMagicLink *schema.MagicLink
	result := ml.db.Where(
		&schema.MagicLink{
			Token: token,
		},
	).Find(
		&dbMagicLink,
	)

	if result.Error != nil {
		log.ErrorLogger.Println(result.Error)
		return dbMagicLink, result.Error
	}
	if result.RowsAffected < 1 {
		log.DebugLogger.Println("no db changes")
		return dbMagicLink, gorm.ErrRecordNotFound
	}

	return dbMagicLink, nil

}

func (ml *MagicLinkAccessOperator) DeleteMagicLink(dbModel *schema.MagicLink) (*schema.MagicLink, error) {
	// var dbmagi
	result := ml.db.Delete(&dbModel)
	if result.Error != nil {
		return dbModel, result.Error
	}
	return dbModel, nil

}
