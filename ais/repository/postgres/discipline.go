package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Дисциплина(
//     id SERIAL,
//     название varchar(150),
//     часы int,
//     CONSTRAINT дисциплина_pk PRIMARY KEY (id)
// );
type Discipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents []*ControlEvent
}

const disciplineIDField = "id"
const disciplineFields = "id,название,часы"
const disciplineTable = "Дисциплина"

func (d *Discipline) toModel(controlEventRef *models.ControlEvent) *models.Discipline {
	discipline := &models.Discipline{
		ID:    d.ID,
		Name:  d.Name,
		Hours: d.Hours,
	}

	controlEvents := make([]*models.ControlEvent, 0)

	for _, controlEvent := range d.ControlEvents {
		if controlEventRef != nil {
			if controlEvent.ID == controlEventRef.ID {
				controlEvents = append(controlEvents, controlEventRef)
			} else {
				controlEvents = append(controlEvents, controlEvent.toModel(discipline, nil, nil))
			}
		} else {
			controlEvents = append(controlEvents, controlEvent.toModel(discipline, nil, nil))
		}
	}

	discipline.ControlEvents = controlEvents

	return discipline
}

func NewPostgresDiscipline(scannable postgres.Scannable) (*Discipline, error) {
	discipline := &Discipline{}

	err := scannable.Scan(&discipline.ID, &discipline.Name, &discipline.Hours)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrDisciplineNotFound
		}
		return nil, err
	}

	return discipline, nil
}

func (d *Discipline) Associate(ctx context.Context, r DBAisRepository, controlEventRef *ControlEvent) error {
	controlEventsRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(
			controlEventTable, controlEventFields, controlEventDisciplineFK,
			disciplineTable, disciplineIDField, disciplineIDField,
		),
		d.ID,
	)

	if err != nil {
		return err
	}

	controlEvents := make([]*ControlEvent, 0)
	for controlEventsRows.Next() {
		controlEvent, err := NewPostgresControlEvent(controlEventsRows)

		if err != nil {
			return err
		}

		if controlEventRef == nil {
			controlEvent.Associate(ctx, r, d, nil, nil)
		} else {
			if controlEventRef.ID == controlEvent.ID {
				controlEvent = controlEventRef
			} else {
				controlEvent.Associate(ctx, r, d, nil, nil)
			}
		}

		controlEvents = append(controlEvents, controlEvent)
	}

	d.ControlEvents = controlEvents

	return nil
}

func makeDisciplineModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Discipline, error) {
	discipline, err := NewPostgresDiscipline(scannable)

	if err != nil {
		return nil, err
	}

	err = discipline.Associate(ctx, r, nil)

	if err != nil {
		return nil, err
	}

	return discipline.toModel(nil), nil
}

func (r DBAisRepository) GetDiscipline(ctx context.Context, disciplineID int) (*models.Discipline, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", disciplineFields, disciplineTable), disciplineID)
	return makeDisciplineModel(ctx, r, row)
}

func (r DBAisRepository) GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error) {
	errValue := []*models.Discipline{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", disciplineFields, disciplineTable))

	if err != nil {
		return errValue, err
	}

	disciplines := make([]*models.Discipline, 0)
	for rows.Next() {
		discipline, err := makeDisciplineModel(ctx, r, rows)

		if err != nil {
			return errValue, err
		}

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}
