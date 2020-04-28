package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type ControlEventType struct {
	ID  int
	Def string
}

func toPostgresControlEventType(c *models.ControlEventType) *ControlEventType {
	return &ControlEventType{
		c.ID,
		c.Def,
	}
}

func toModelControlEventType(c *ControlEventType) *models.ControlEventType {
	return &models.ControlEventType{
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

const getControlEventTypeQuery = `SELECT * FROM ТипКонтрольногоМероприятия WHERE id = $1`

func (r DBAisRepository) GetControlEventType(ctx context.Context, controlEventTypeID int) (*models.ControlEventType, error) {
	row := r.db.QueryRowContext(ctx, getControlEventTypeQuery, controlEventTypeID)

	controlEventType := new(ControlEventType)
	err := row.Scan(&controlEventType.ID, &controlEventType.Def)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrControlEventTypeNotFound
		}
		return nil, err
	}

	return toModelControlEventType(controlEventType), nil
}

const getAllControlEventTypesQuery = `SELECT * FROM ТипКонтрольногоМероприятия`

func (r DBAisRepository) GetAllControlEventTypes(ctx context.Context) ([]*models.ControlEventType, error) {
	rows, err := r.db.QueryContext(ctx, getAllControlEventTypesQuery)
	controlEventTypes := make([]*models.ControlEventType, 0)

	if err != nil {
		return controlEventTypes, err
	}

	for rows.Next() {
		controlEventType := new(ControlEventType)
		if err := rows.Scan(&controlEventType.ID, &controlEventType.Def); err != nil {
			return []*models.ControlEventType{}, err
		}
		controlEventTypes = append(controlEventTypes, toModelControlEventType(controlEventType))
	}

	return controlEventTypes, nil
}
