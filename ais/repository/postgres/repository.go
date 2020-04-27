package postgres

import (
	"database/sql"
)

type AisRepository struct {
	db *sql.DB
}

func NewAisRepository(db *sql.DB) *AisRepository {
	return &AisRepository{
		db: db,
	}
}
