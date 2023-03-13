package migration

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	user = `
	CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
		role INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP,
	);
	`
	email = `
	CREATE TABLE IF NOT EXISTS emails (
		id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
		user_id uuid NOT NULL,
		email TEXT NOT NULL UNIQUE,
		verified BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);
	`
	magicLink = `
	CREATE TABLE IF NOT EXISTS magic_links (
		id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
		user_id uuid NOT NULL,
		email_id uuid NOT NULL,
		token uuid NOT NULL,
		link_type INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		expiry_time TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
		FOREIGN KEY (email_id) REFERENCES emails (id)
	);
	`
	credential = `
	CREATE TABLE IF NOT EXISTS credentials (
		id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
		user_id uuid NOT NULL,
		auth_type INTEGER NOT NULL,
		provider_type INTEGER NOT NULL,
		identifier TEXT NOT NULL,
		password_hash TEXT,
		password_algo TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users (id),
	);
	`
)

func InitDB(dataSourceName string) (*ModelDao, error) {
	var err error

	db, err := sqlx.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to Open sqlite3 DB")
	}
	db.SetMaxOpenConns(10)

	table_schema := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS invites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			token TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			role TEXT NOT NULL,
			org_id TEXT NOT NULL,
			FOREIGN KEY(org_id) REFERENCES organizations(id)
		);
		CREATE TABLE IF NOT EXISTS organizations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			is_anonymous INTEGER NOT NULL DEFAULT 0 CHECK(is_anonymous IN (0,1)),
			has_opted_updates INTEGER NOT NULL DEFAULT 1 CHECK(has_opted_updates IN (0,1))
		);
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			profile_picture_url TEXT,
			group_id TEXT NOT NULL,
			org_id TEXT NOT NULL,
			FOREIGN KEY(group_id) REFERENCES groups(id),
			FOREIGN KEY(org_id) REFERENCES organizations(id)
		);
		CREATE TABLE IF NOT EXISTS groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS reset_password_request (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			token TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS user_flags (
			user_id TEXT PRIMARY KEY,
			flags TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`

	_, err = db.Exec(table_schema)
	if err != nil {
		return nil, fmt.Errorf("Error in creating tables: %v", err.Error())
	}

	mds := &ModelDaoSqlite{db: db}

	ctx := context.Background()
	if err := mds.initializeOrgPreferences(ctx); err != nil {
		return nil, err
	}
	if err := mds.initializeRBAC(ctx); err != nil {
		return nil, err
	}

	return mds, nil
}
