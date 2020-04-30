package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Группа(
//     id SERIAL,
//     id_кафедры int NOT NULL references Кафедра(id) ON DELETE CASCADE,
//     id_семестра int NOT NULL references Семестр(id) ON DELETE CASCADE,
//     CONSTRAINT группа_pk PRIMARY KEY (id)
// );

type Group struct {
	ID     int
	Number int

	Cathedra  *Cathedra
	Students  []*Student
	Semesters []*Semester
}

const groupTable = "Группа"
const groupIDField = "id"
const groupFields = "id,номер"
const groupCathedraFK = "id_кафедры"

const groupSemesterMtM = "Группа__Семестр"
const groupSemesterMTMGroupKey = "id_группы"
const groupSemesterMTMSemesterKey = "id_семестра"

// const studentResidenceFK = "id_места_жительства"

// func toModelGroup(r DBAisRepository, ctx context.Context, g *Group, callerStudent *models.Student) (*models.Group, error) {
func (g *Group) toModel(
	studentRef *models.Student, cathedraRef *models.Cathedra,
	semesterRef *models.Semester,
) *models.Group {
	group := &models.Group{
		ID:       g.ID,
		Number:   g.Number,
		Cathedra: cathedraRef,
	}

	students := make([]*models.Student, 0)
	for _, student := range g.Students {
		if studentRef != nil {
			if student.ID == studentRef.ID {
				students = append(students, studentRef)
			} else {
				students = append(students, student.toModel(nil, group, nil, nil))
			}
		} else {
			students = append(students, student.toModel(nil, group, nil, nil))
		}
	}
	group.Students = students

	if group.Cathedra == nil {
		group.Cathedra = g.Cathedra.toModel(group)
	}

	semesters := make([]*models.Semester, 0)
	for _, semester := range g.Semesters {
		if semesterRef != nil {
			if semester.ID == semesterRef.ID {
				semesters = append(semesters, semesterRef)
			} else {
				semesters = append(semesters, semester.toModel(group, nil))
			}
		} else {
			semesters = append(semesters, semester.toModel(group, nil))
		}
	}
	group.Semesters = semesters

	return group
}

func NewPostgresGroup(scannable postgres.Scannable) (*Group, error) {
	group := &Group{}

	err := scannable.Scan(&group.ID, &group.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrGroupNotFound
		}
		return nil, err
	}

	return group, nil
}

func (g *Group) Associate(
	ctx context.Context, r DBAisRepository,
	studentRef *Student, cathedraRef *Cathedra, semesterRef *Semester,
) error {
	studentsRow, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(studentTable, studentFields, studentGroupFK, groupTable, groupIDField, groupIDField),
		g.ID,
	)

	if err != nil {
		return err
	}

	students := make([]*Student, 0)
	for studentsRow.Next() {
		student, err := NewPostgresStudent(studentsRow)

		if err != nil {
			return err
		}

		if studentRef == nil {
			student.Associate(ctx, r, nil, g, nil, nil)
		} else {
			if studentRef.ID == student.ID {
				student = studentRef
			} else {
				student.Associate(ctx, r, nil, g, nil, nil)
			}
		}

		students = append(students, student)
	}

	g.Students = students

	if cathedraRef == nil {
		cathedraRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(cathedraTable, cathedraFields, "id", groupTable, groupCathedraFK, "id"),
			g.ID,
		)

		cathedra, err := NewPostgresCathedra(cathedraRow)

		if err != nil {
			return err
		}

		cathedra.Associate(ctx, r, g)
		g.Cathedra = cathedra
	} else {
		g.Cathedra = cathedraRef
	}

	semesterRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeManyToManyJoinQuery(
			semesterTable, semesterFields, semesterIDField, groupSemesterMTMSemesterKey,
			groupTable, groupIDField, groupSemesterMTMGroupKey,
			groupSemesterMtM,
		),
		g.ID,
	)

	if err != nil {
		return err
	}

	semesters := make([]*Semester, 0)
	for semesterRows.Next() {
		semester, err := NewPostgresSemester(semesterRows)

		if err != nil {
			return err
		}

		if semesterRef == nil {
			semester.Associate(ctx, r, nil, g)
		} else {
			if semesterRef.ID == semester.ID {
				semester = semesterRef
			} else {
				semester.Associate(ctx, r, nil, g)
			}
		}

		semesters = append(semesters, semester)
	}
	g.Semesters = semesters

	return nil
}

func makeGroupModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Group, error) {
	group, err := NewPostgresGroup(scannable)

	if err != nil {
		return nil, err
	}

	err = group.Associate(ctx, r, nil, nil, nil)

	if err != nil {
		return nil, err
	}

	return group.toModel(nil, nil, nil), nil
}

func (r DBAisRepository) GetGroup(ctx context.Context, groupID int) (*models.Group, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", groupFields, groupTable), groupID)
	return makeGroupModel(ctx, r, row)
}

func (r DBAisRepository) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	errValue := []*models.Group{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", groupFields, groupTable))

	if err != nil {
		return errValue, err
	}

	groups := []*models.Group{}
	for rows.Next() {
		group, err := makeGroupModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		groups = append(groups, group)
	}

	return groups, nil
}
