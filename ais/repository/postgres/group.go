package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

// CREATE TABLE Группа(
//     id SERIAL,
//     id_кафедры int NOT NULL references Кафедра(id) ON DELETE CASCADE,
//     id_семестра int NOT NULL references Семестр(id) ON DELETE CASCADE,
//     CONSTRAINT группа_pk PRIMARY KEY (id)
// );

type Group struct {
	ID          int
	CathedraID  int
	Number      int
	SemesterIDs []int
	StudentIDs  []int
}

func toPostgresGroup(g *models.Group) *Group {
	return nil
}

func toModelGroup(r DBAisRepository, ctx context.Context, g *Group, callerStudent *models.Student) (*models.Group, error) {
	semesterModels := make([]*models.Semester, 0, len(g.SemesterIDs))

	for _, semesterID := range g.SemesterIDs {
		semesterModel, err := r.GetSemester(ctx, semesterID)
		if err != nil {
			return nil, err
		}

		semesterModels = append(semesterModels, semesterModel)
	}

	cathedraModel, err := r.GetCathedra(ctx, g.CathedraID)
	if err != nil {
		return nil, err
	}

	studentModels := make([]*models.Student, 0, len(g.StudentIDs))

	for _, studentID := range g.StudentIDs {
		if callerStudent != nil {
			if callerStudent.ID == studentID {
				continue
			}
		}

		studentModel, err := r.GetStudent(ctx, studentID)
		if err != nil {
			return nil, err
		}

		studentModels = append(studentModels, studentModel)
	}

	return &models.Group{
		ID:        g.ID,
		Cathedra:  cathedraModel,
		Students:  studentModels,
		Semesters: semesterModels,
		Number:    g.Number,
	}, nil

}

const getStudentIDsInGroupQuery = `
SELECT с.id FROM
	Группа as г
	JOIN Студент as с ON г.id = с.id_группы
	WHERE г.id = $1`

func (g *Group) Fill(ctx context.Context, db *sql.DB, sc Scannable) error {
	err := g.hydrate(sc)

	if err != nil {
		if err == sql.ErrNoRows {
			return ais.ErrStudentNotFound
		}
		return err
	}

	idRows, err := db.QueryContext(ctx, getStudentIDsInGroupQuery, g.ID)

	if err != nil {
		return err
	}

	for idRows.Next() {
		var studentID int
		idRows.Scan(&studentID)
		g.StudentIDs = append(g.StudentIDs, studentID)
	}

	return nil
}

const getGroupQuery = `SELECT id, id_кафедры, номер FROM Группа WHERE id = $1`

func (g *Group) hydrate(sc Scannable) error {
	return sc.Scan(&g.ID, &g.CathedraID, &g.Number)
}

func (r DBAisRepository) GetGroupRecursive(ctx context.Context, groupID int, caller *models.Student) (*models.Group, error) {
	row := r.db.QueryRowContext(ctx, getGroupQuery, groupID)

	group := new(Group)
	err := group.Fill(ctx, r.db, row)

	if err != nil {
		return nil, err
	}

	groupModel, err := toModelGroup(r, ctx, group, caller)

	if err != nil {
		return nil, err
	}

	return groupModel, nil
}
func (r DBAisRepository) GetGroup(ctx context.Context, groupID int) (*models.Group, error) {
	row := r.db.QueryRowContext(ctx, getGroupQuery, groupID)

	group := new(Group)
	err := group.Fill(ctx, r.db, row)

	if err != nil {
		return nil, err
	}

	groupModel, err := toModelGroup(r, ctx, group, nil)

	if err != nil {
		return nil, err
	}

	return groupModel, nil
}

const getAllGroupsQuery = `SELECT id, id_кафедры, номер FROM Группа`

func (r DBAisRepository) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	rows, err := r.db.QueryContext(ctx, getAllGroupsQuery)
	groups := make([]*models.Group, 0)

	if err != nil {
		return groups, err
	}

	for rows.Next() {
		group := new(Group)
		err := group.Fill(ctx, r.db, rows)

		if err != nil {
			return nil, err
		}

		groupModel, err := toModelGroup(r, ctx, group, nil)
		if err != nil {
			return nil, err
		}

		groups = append(groups, groupModel)

	}

	return groups, nil
}
