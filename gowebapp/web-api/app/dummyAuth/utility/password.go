package utility

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// func setPassword(userId uuid.UUID, password string, w http.ResponseWriter, db *gorm.DB) {
// 	var dbEmail models.Email
// 	// var dbUser models.User
// 	var dbCredentials models.Credentials
// 	var hashPass []byte
// 	var err error

// 	dbEmail, err = dao.FindPrimaryEmailbyUserId(userId, db)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	hashPass, err = Hash(password)
// 	if err != nil {
// 		log.Println(err.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	result := db.Where(&models.Credentials{UserId: userId, ProviderType: uint8(constants.DBPassword), AuthType: uint8(constants.Custom)}).First(&dbCredentials)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		log.Println(result.Error)
// 		// create new entry credentials

// 		dbCredentials, err = dao.CreateCredentials(userId, string(hashPass), "bcrypt", dbEmail.Email, constants.Custom, constants.DBPassword, db)
// 		log.Println("new credentials record created!!!!")
// 		responsePayload := SignInResponseBody{IsSuccessful: true, UserId: userId.String(), BearerToken: "bearer-willcomesoon"}
// 		json.NewEncoder(w).Encode(&responsePayload)
// 		return

// 	}

// 	//  update credentials
// 	algo := "bcrypt"
// 	hashPassString := string(hashPass)
// 	result = db.Model(&dbCredentials).Updates(models.Credentials{PasswordHash: &hashPassString, PasswordAlgo: &algo, Identifier: dbEmail.Email})
// 	if result.Error != nil {
// 		log.Println(result.Error.Error())
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	log.Println("credentials updated!!!!")
// 	responsePayload := SignInResponseBody{IsSuccessful: true, UserId: userId.String(), BearerToken: "bearer-willcomesoon"}
// 	json.NewEncoder(w).Encode(&responsePayload)
// 	return

// }
