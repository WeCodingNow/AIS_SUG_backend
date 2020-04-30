package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Семестр(
//     id SERIAL,
//     номер int NOT NULL,
//     начало date NOT NULL,
//     конец date,
//     CONSTRAINT семестр_pk PRIMARY KEY (id)
// );
type Semester struct {
	ID        int
	Number    int
	Beginning sql.NullTime
	End       sql.NullTime

	Groups        []*Group
	ControlEvents []*ControlEvent
}

const semesterIDField = "id"
const semesterFields = "id,номер,начало,конец"
const semesterTable = "Семестр"

func (s *Semester) toModel(
	groupRef *models.Group, controlEventRef *models.ControlEvent,
) *models.Semester {
	semester := &models.Semester{
		ID:        s.ID,
		Number:    s.Number,
		Beginning: s.Beginning.Time,
	}

	if s.End.Valid {
		semester.End = new(time.Time)
		*semester.End = s.End.Time
	}

	groups := make([]*models.Group, 0)
	for _, group := range s.Groups {
		if groupRef != nil {
			if group.ID == groupRef.ID {
				groups = append(groups, groupRef)
			} else {
				groups = append(groups, group.toModel(nil, nil, semester))
			}
		} else {
			groups = append(groups, group.toModel(nil, nil, semester))
		}
	}
	semester.Groups = groups

	// log.Println("THIS RUNS")

	// controlEvents := make([]*models.ControlEvent, 0)
	// for _, controlEvent := range s.ControlEvents {
	// 	if controlEventRef != nil {
	// 		if controlEvent.ID == controlEventRef.ID {
	// 			controlEvents = append(controlEvents, controlEventRef)
	// 		} else {
	// 			controlEvents = append(controlEvents, controlEvent.toModel(nil, semester, nil))
	// 		}
	// 	} else {
	// 		controlEvents = append(controlEvents, controlEvent.toModel(nil, semester, nil))
	// 	}
	// }
	// semester.ControlEvents = controlEvents

	return semester
}

func NewPostgresSemester(scannable postgres.Scannable) (*Semester, error) {
	semester := &Semester{}

	err := scannable.Scan(&semester.ID, &semester.Number, &semester.Beginning, &semester.End)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrSemesterNotFound
		}
		return nil, err
	}

	return semester, nil
}

func (s *Semester) Associate(
	ctx context.Context, r DBAisRepository,
	controlEventRef *ControlEvent, groupRef *Group,
) error {
	controlEventsRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(
			controlEventTable, controlEventFields, controlEventSemesterFK,
			semesterTable, semesterIDField, semesterIDField,
		),
		s.ID,
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
			controlEvent.Associate(ctx, r, nil, s, nil)
		} else {
			if controlEventRef.ID == controlEvent.ID {
				controlEvent = controlEventRef
			} else {
				controlEvent.Associate(ctx, r, nil, s, nil)
			}
		}

		controlEvents = append(controlEvents, controlEvent)
	}
	s.ControlEvents = controlEvents

	groupRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeManyToManyJoinQuery(
			groupTable, groupFields, groupIDField, groupSemesterMTMGroupKey,
			semesterTable, semesterIDField, groupSemesterMTMSemesterKey,
			groupSemesterMtM,
		),
		s.ID,
	)

	if err != nil {
		return err
	}

	groups := make([]*Group, 0)
	for groupRows.Next() {
		group, err := NewPostgresGroup(groupRows)

		if err != nil {
			return err
		}

		if groupRef == nil {
			group.Associate(ctx, r, nil, nil, s)
		} else {
			if groupRef.ID == group.ID {
				group = groupRef
			} else {
				group.Associate(ctx, r, nil, nil, s)
			}
		}

		groups = append(groups, group)
	}
	s.Groups = groups

	return nil
}

func makeSemesterModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Semester, error) {
	semester, err := NewPostgresSemester(scannable)

	if err != nil {
		return nil, err
	}

	err = semester.Associate(ctx, r, nil, nil)

	if err != nil {
		return nil, err
	}

	return semester.toModel(nil, nil), nil
}

func (r DBAisRepository) GetSemester(ctx context.Context, semesterID int) (*models.Semester, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", semesterFields, semesterTable), semesterID)
	return makeSemesterModel(ctx, r, row)
}

func (r DBAisRepository) GetAllSemesters(ctx context.Context) ([]*models.Semester, error) {
	errValue := []*models.Semester{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", semesterFields, semesterTable))

	if err != nil {
		return errValue, err
	}

	semesters := make([]*models.Semester, 0)
	for rows.Next() {
		semester, err := makeSemesterModel(ctx, r, rows)

		if err != nil {
			return errValue, err
		}

		semesters = append(semesters, semester)
	}

	return semesters, nil
}
