package dao

import (
	"web-api/app/dummyAuth/schema"
	log "web-api/app/utility/logger"

	// "web-api/app/utility/email"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (pda *ProfileDataAccess) OnboardingTransaction(
	dbProfile *schema.Profile,
	profileUdateModel *schema.ProfileBase,
	emailId uuid.UUID,
) (*schema.Profile, *schema.AccessRequest, error) {
	var dbAccessRequest *schema.AccessRequest

	err := pda.db.Transaction(func(tx *gorm.DB) error {

		// // Create Password
		// algoPass := constants.Bcrypt.String()
		// dbCredential := schema.Credential{
		// 	UserId:       dbProfile.UserId,
		// 	AuthType:     uint8(constants.Custom),
		// 	ProviderType: uint8(constants.DBPassword),
		// 	PasswordHash: hashPassword,
		// 	Identifier:   emailId.String(),
		// 	PasswordAlgo: &algoPass,
		// }

		// result := tx.Create(&dbCredential)
		// if result.Error != nil {
		// 	return result.Error
		// }

		// Create Access Request
		dbAccessRequest = &schema.AccessRequest{ProfileId: dbProfile.ID, Approved: false, EmailId: emailId}
		result := tx.Create(&dbAccessRequest)
		if result.Error != nil {
			return result.Error
		}

		// Update Profile
		result = tx.Model(&dbProfile).Updates(schema.Profile{ProfileBase: *profileUdateModel, Onboarded: true})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		log.ErrorLogger.Println("errorn in transaction")
		return dbProfile, dbAccessRequest, err
	}

	// email.SendEmailviaAws(dbEmail.Email, nil, constants.SuccessMessage, nil, constants.Onboarding, pda.cfg)

	return dbProfile, dbAccessRequest, nil
}
