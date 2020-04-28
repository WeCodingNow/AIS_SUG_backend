package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type ContactType struct {
	ID  int
	Def string
}

// CREATE TABLE Кафедра(
//     id SERIAL,
//     название varchar(100) NOT NULL UNIQUE,
//     короткое_название varchar(10) NOT NULL UNIQUE,
//     CONSTRAINT кафедра_pk PRIMARY KEY (id)
// );

func toPostgresContactType(c *models.ContactType) *ContactType {
	return &ContactType{
		c.ID,
		c.Def,
	}
}

func toModelContactType(c *ContactType) *models.ContactType {
	return &models.ContactType{
		c.ID,
		c.Def,
	}
}

// const createCathedraQuery = `INSERT INTO Кафедра(название, короткое_название) VALUES ( $1, $2 )`

// func (r AisRepository) CreateCathedra(ctx context.Context, name, shortName string) error {
// 	_, err := r.db.ExecContext(ctx, createCathedraQuery,
// 		name, shortName,
// 	)

// 	return err
// }

const getContactTypeQuery = `SELECT * FROM ТипКонтакта WHERE id = $1`

func (r DBAisRepository) GetContactType(ctx context.Context, contactTypeID int) (*models.ContactType, error) {
	row := r.db.QueryRowContext(ctx, getContactTypeQuery, contactTypeID)

	contactType := new(ContactType)
	err := row.Scan(&contactType.ID, &contactType.Def)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrContactTypeNotFound
		}
		return nil, err
	}

	return toModelContactType(contactType), nil
}

const getAllContactTypesQuery = `SELECT * FROM ТипКонтакта`

func (r DBAisRepository) GetAllContactTypes(ctx context.Context) ([]*models.ContactType, error) {
	rows, err := r.db.QueryContext(ctx, getAllContactTypesQuery)
	contactTypes := make([]*models.ContactType, 0)

	if err != nil {
		return contactTypes, err
	}

	for rows.Next() {
		contactType := new(ContactType)
		if err := rows.Scan(&contactType.ID, &contactType.Def); err != nil {
			return []*models.ContactType{}, err
		}
		contactTypes = append(contactTypes, toModelContactType(contactType))
	}

	return contactTypes, nil
}
