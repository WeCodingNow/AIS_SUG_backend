package postgres

import (
	"database/sql"
)

type AisAuthRepository struct {
	db *sql.DB
}

func NewAisAuthUserRepository(db *sql.DB) *AisAuthRepository {
	return &AisAuthRepository{
		db: db,
	}
}
