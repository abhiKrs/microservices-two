package database

import (
	// "context"
	"fmt"
	// "log"
	"time"

	"web-api/app/config"
	"web-api/app/dummyAuth/schema"
	sourceSchema "web-api/app/source/schema"
	teamSchema "web-api/app/team/schema"
	logs "web-api/app/utility/logger"
	viewSchema "web-api/app/view/schema"

	// "web-api/test/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(opts config.PGDatabase) *gorm.DB {
	var err error

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Calcutta",
		opts.HOST,
		opts.USER,
		opts.PASSWORD,
		opts.DB,
		opts.PORT,
	)
	// defer handleNoDB()
	logs.DebugLogger.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

	if err != nil {

		logs.ErrorLogger.Println(err)
		logs.ErrorLogger.Println("database not found")
		panic(err)
	}
	// fmt.Println("? Connected Successfully to the Database")
	// // a := string("uuid-ossp")
	// query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	// log.Println(query)
	// db = db.Exec(query)

	db.AutoMigrate(
		&schema.User{},
		&schema.Email{},
		&schema.MagicLink{},
		&schema.Credential{},
		&schema.Profile{},
		&schema.AccessRequest{},
		&sourceSchema.Source{},
		&teamSchema.Team{},
		&viewSchema.View{},
	)

	sqlDB, err := db.DB()
	if err != nil {
		// control error
		logs.DebugLogger.Println(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return db
}
