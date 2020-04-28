package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type ControlEvent struct {
	ID                 int
	ControlEventTypeID int
	DisciplineID       int
	SemesterID         int
	Date               time.Time
}

func toPostgresControlEvent(c *models.ControlEvent) *ControlEvent {
	return nil
}

func toModelControlEvent(r AisRepository, ctx context.Context, c *ControlEvent) (*models.ControlEvent, error) {
	controlEventType, err := r.GetControlEventType(ctx, c.ControlEventTypeID)
	if err != nil {
		return nil, err
	}

	discipline, err := r.GetDiscipline(ctx, c.DisciplineID)
	if err != nil {
		return nil, err
	}
	semester, err := r.GetSemester(ctx, c.SemesterID)
	if err != nil {
		return nil, err
	}

	return &models.ControlEvent{
		ID:               c.ID,
		Date:             c.Date,
		ControlEventType: controlEventType,
		Discipline:       discipline,
		Semester:         semester,
	}, nil
}

const getControlEventQuery = `SELECT id, id_типа, id_дисциплины, id_семестра, дата_проведения FROM КонтрольноеМероприятие WHERE id = $1`

func (r AisRepository) GetControlEvent(ctx context.Context, controlEventID int) (*models.ControlEvent, error) {
	row := r.db.QueryRowContext(ctx, getControlEventQuery, controlEventID)

	controlEvent := new(ControlEvent)
	err := row.Scan(&controlEvent.ID, &controlEvent.ControlEventTypeID, &controlEvent.DisciplineID, &controlEvent.SemesterID, &controlEvent.Date)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrControlEventNotFound
		}
		return nil, err
	}

	return toModelControlEvent(r, ctx, controlEvent)
}

const getAllControlEventsQuery = `SELECT id, id_типа, id_дисциплины, id_семестра, дата_проведения FROM КонтрольноеМероприятие`

func (r AisRepository) GetAllControlEvents(ctx context.Context) ([]*models.ControlEvent, error) {
	rows, err := r.db.QueryContext(ctx, getAllControlEventsQuery)
	controlEvents := make([]*models.ControlEvent, 0)

	if err != nil {
		return controlEvents, err
	}

	for rows.Next() {
		controlEvent := new(ControlEvent)
		if err := rows.Scan(&controlEvent.ID, &controlEvent.ControlEventTypeID, &controlEvent.DisciplineID, &controlEvent.SemesterID, &controlEvent.Date); err != nil {
			return []*models.ControlEvent{}, err
		}
		controlEventModel, err := toModelControlEvent(r, ctx, controlEvent)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ais.ErrControlEventNotFound
			}
			return nil, err
		}

		controlEvents = append(controlEvents, controlEventModel)
	}

	return controlEvents, nil
}
