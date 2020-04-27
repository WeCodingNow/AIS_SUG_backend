package postgres

import "github.com/WeCodingNow/AIS_SUG_backend/models"

// CREATE TABLE Группа(
//     id SERIAL,
//     id_кафедры int NOT NULL references Кафедра(id) ON DELETE CASCADE,
//     id_семестра int NOT NULL references Семестр(id) ON DELETE CASCADE,
//     CONSTRAINT группа_pk PRIMARY KEY (id)
// );

type Group struct {
	ID         int
	SemesterID int
	CathedraID int
}

func toPostgresGroup(g *models.Group) *Group {
	return nil
}

func toModelGroup(g *Group) *models.Group {
	return nil
}
