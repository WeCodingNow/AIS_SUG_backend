package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE ТипКонтакта(
//     id SERIAL,
//     обозначение varchar(100) NOT NULL,
//     CONSTRAINT тип_контакта_pk PRIMARY KEY (id)
// );
type ContactType struct {
	ID  int
	Def string
}

const contactTypeIDField = "id"
const contactTypeFields = "id,обозначение"
const contactTypeTable = "ТипКонтакта"

func (c *ContactType) toModel() *models.ContactType {
	return &models.ContactType{
		ID:  c.ID,
		Def: c.Def,
	}
}

func NewPostgresContactType(scannable postgres.Scannable) (*ContactType, error) {
	contactType := &ContactType{}

	err := scannable.Scan(&contactType.ID, &contactType.Def)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrContactTypeNotFound
		}
		return nil, err
	}

	return contactType, nil
}

func (r DBAisRepository) GetContactType(ctx context.Context, contactTypeID int) (*models.ContactType, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", contactTypeFields, contactTypeTable), contactTypeID)

	contactType, err := NewPostgresContactType(row)

	if err != nil {
		return nil, err
	}

	return contactType.toModel(), nil
}

func (r DBAisRepository) GetAllContactTypes(ctx context.Context) ([]*models.ContactType, error) {
	errValue := []*models.ContactType{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", contactTypeFields, contactTypeTable))

	if err != nil {
		return errValue, err
	}

	contactTypes := make([]*models.ContactType, 0)
	for rows.Next() {
		contactType, err := NewPostgresContactType(rows)
		if err != nil {
			return errValue, err
		}
		contactTypes = append(contactTypes, contactType.toModel())
	}

	return contactTypes, nil
}
