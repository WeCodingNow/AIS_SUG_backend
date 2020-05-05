package postgres

import (
	"database/sql"
)

type DBAisRepository struct {
	db *sql.DB
}

func NewDBAisRepository(db *sql.DB) *DBAisRepository {
	return &DBAisRepository{
		db: db,
	}
}
