package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE КонтрольноеМероприятие(
//     id SERIAL,
//     id_типа int NOT NULL references ТипКонтрольногоМероприятия(id) ON DELETE CASCADE,
//     id_дисциплины int NOT NULL references Дисциплина(id) ON DELETE CASCADE,
//     id_семестра int NOT NULL references Семестр(id) ON DELETE CASCADE,
//     дата_проведения date NOT NULL,
//     CONSTRAINT контрольное_мероприятие_pk PRIMARY KEY (id)
// );
type ControlEvent struct {
	ID   int
	Date sql.NullTime

	ControlEventType *ControlEventType
	Discipline       *Discipline
	// Semester         *Semester
	Marks []*Mark
}

const controlEventTable = "КонтрольноеМероприятие"
const controlEventIDField = "id"
const controlEventFields = "id,дата_проведения"
const controlEventControlEventTypeFK = "id_типа"
const controlEventDisciplineFK = "id_дисциплины"
const controlEventSemesterFK = "id_семестра"

func (ce *ControlEvent) toModel(
	disciplineRef *models.Discipline, semesterRef *models.Semester,
	markRef *models.Mark,
) *models.ControlEvent {
	controlEvent := &models.ControlEvent{
		ID:   ce.ID,
		Date: ce.Date.Time,

		ControlEventType: ce.ControlEventType.toModel(),
		Discipline:       disciplineRef,
		// Semester:         semesterRef,
	}

	if controlEvent.Discipline == nil {
		controlEvent.Discipline = ce.Discipline.toModel(controlEvent)
	}

	// if controlEvent.Semester == nil {
	// 	controlEvent.Semester = ce.Semester.toModel(nil, controlEvent)
	// }

	marks := make([]*models.Mark, 0)
	for _, mark := range ce.Marks {
		if markRef != nil {
			if mark.ID == markRef.ID {
				marks = append(marks, markRef)
			} else {
				marks = append(marks, mark.toModel(controlEvent, nil))
			}
		} else {
			marks = append(marks, mark.toModel(controlEvent, nil))
		}
	}
	controlEvent.Marks = marks

	return controlEvent
}

func NewPostgresControlEvent(scannable postgres.Scannable) (*ControlEvent, error) {
	controlEvent := &ControlEvent{}

	err := scannable.Scan(&controlEvent.ID, &controlEvent.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrControlEventNotFound
		}
		return nil, err
	}

	return controlEvent, nil
}

func (ce *ControlEvent) Associate(
	ctx context.Context, r DBAisRepository,
	disciplineRef *Discipline, semesterRef *Semester, markRef *Mark,
) error {
	controlEventTypeRow := r.db.QueryRowContext(
		ctx,
		postgres.MakeJoinQuery(
			controlEventTypeTable, controlEventTypeFields, controlEventTypeIDField,
			controlEventTable, controlEventControlEventTypeFK, controlEventIDField),
		ce.ID,
	)

	controlEventType, err := NewPostgresControlEventType(controlEventTypeRow)

	if err != nil {
		return err
	}

	ce.ControlEventType = controlEventType

	if disciplineRef == nil {
		disciplineRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(
				disciplineTable, disciplineFields, "id",
				controlEventTable, controlEventDisciplineFK, "id"),
			ce.ID,
		)

		discipline, err := NewPostgresDiscipline(disciplineRow)

		if err != nil {
			return err
		}

		discipline.Associate(ctx, r, ce)

		ce.Discipline = discipline
	} else {
		ce.Discipline = disciplineRef
	}

	// if semesterRef == nil {
	// 	semesterRow := r.db.QueryRowContext(
	// 		ctx,
	// 		postgres.MakeJoinQuery(
	// 			semesterTable, semesterFields, "id",
	// 			controlEventTable, controlEventSemesterFK, "id"),
	// 		ce.ID,
	// 	)

	// 	semester, err := NewPostgresSemester(semesterRow)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	semester.Associate(ctx, r, ce, nil)
	// 	ce.Semester = semester
	// } else {
	// 	ce.Semester = semesterRef
	// }

	markRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(markTable, markFields, markControlEventFK, controlEventTable, controlEventIDField, markIDField),
		ce.ID,
	)

	if err != nil {
		return err
	}

	marks := make([]*Mark, 0)
	for markRows.Next() {
		mark, err := NewPostgresMark(markRows)

		if err != nil {
			return err
		}

		if markRef == nil {
			mark.Associate(ctx, r, ce, nil)
		} else {
			if markRef.ID == mark.ID {
				mark = markRef
			} else {
				mark.Associate(ctx, r, ce, nil)
			}
		}

		marks = append(marks, mark)
	}

	ce.Marks = marks

	return nil
}

func makeControlEventModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.ControlEvent, error) {
	controlEvent, err := NewPostgresControlEvent(scannable)

	if err != nil {
		return nil, err
	}

	err = controlEvent.Associate(ctx, r, nil, nil, nil)

	if err != nil {
		return nil, err
	}

	return controlEvent.toModel(nil, nil, nil), nil
}

func (r DBAisRepository) GetControlEvent(ctx context.Context, controlEventID int) (*models.ControlEvent, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", controlEventFields, controlEventTable), controlEventID)

	return makeControlEventModel(ctx, r, row)
}

func (r DBAisRepository) GetAllControlEvents(ctx context.Context) ([]*models.ControlEvent, error) {
	errValue := []*models.ControlEvent{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", controlEventFields, controlEventTable))

	if err != nil {
		return errValue, err
	}

	controlEvents := []*models.ControlEvent{}
	for rows.Next() {
		controlEvent, err := makeControlEventModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		controlEvents = append(controlEvents, controlEvent)
	}

	return controlEvents, nil
}
