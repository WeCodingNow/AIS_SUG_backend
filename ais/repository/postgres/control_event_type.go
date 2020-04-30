package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE ТипКонтрольногоМероприятия(
//     id SERIAL,
//     обозначение varchar(50),
//     CONSTRAINT тип_контрольного_мероприятие_pk PRIMARY KEY (id)
// );
type ControlEventType struct {
	ID  int
	Def string
}

const controlEventTypeIDField = "id"
const controlEventTypeFields = "id,обозначение"
const controlEventTypeTable = "ТипКонтрольногоМероприятия"

func (c *ControlEventType) toModel() *models.ControlEventType {
	return &models.ControlEventType{
		ID:  c.ID,
		Def: c.Def,
	}
}

func NewPostgresControlEventType(scannable postgres.Scannable) (*ControlEventType, error) {
	contactType := &ControlEventType{}

	err := scannable.Scan(&contactType.ID, &contactType.Def)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrControlEventTypeNotFound
		}
		return nil, err
	}

	return contactType, nil
}

func (r DBAisRepository) GetControlEventType(ctx context.Context, controlEventTypeID int) (*models.ControlEventType, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", controlEventTypeFields, controlEventTypeTable), controlEventTypeID)

	contactType, err := NewPostgresControlEventType(row)

	if err != nil {
		return nil, err
	}

	return contactType.toModel(), nil
}

func (r DBAisRepository) GetAllControlEventTypes(ctx context.Context) ([]*models.ControlEventType, error) {
	errValue := []*models.ControlEventType{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", controlEventTypeFields, controlEventTypeTable))

	if err != nil {
		return errValue, err
	}

	controlEventTypes := make([]*models.ControlEventType, 0)
	for rows.Next() {
		controlEventType, err := NewPostgresControlEventType(rows)
		if err != nil {
			return errValue, err
		}
		controlEventTypes = append(controlEventTypes, controlEventType.toModel())
	}

	return controlEventTypes, nil
}
