package postgres

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
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
	// Cathedra *Cathedra
	Students []*Student
	// Semesters []*Semester
}

func toPostgresGroup(g *models.Group) *Group {
	return nil
}

// func toModelGroup(r DBAisRepository, ctx context.Context, g *Group, callerStudent *models.Student) (*models.Group, error) {
func (g *Group) toModelGroup() (*models.Group, error) {
	return nil, nil
}

func (r DBAisRepository) GetGroup(ctx context.Context, groupID int) (*models.Group, error) {
	return nil, nil
}

const getAllGroupsQuery = `SELECT id, id_кафедры, номер FROM Группа`

func (r DBAisRepository) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	return []*models.Group{}, nil
}
