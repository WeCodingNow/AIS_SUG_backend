package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE Задолженность(
//     id SERIAL,
//     id_дисциплины int NOT NULL references Дисциплина(id) ON DELETE CASCADE,
// 	   id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     описание varchar(200),
//     ликвидирована boolean NOT NULL DEFAULT false,
//     CONSTRAINT задолженность_pk PRIMARY KEY (id)
// );

type repoBacklog struct {
	ID          int
	Description string
	Done        bool

	Discipline *repoDiscipline
	Student    *repoStudent

	model *models.Backlog
}

func NewRepoBacklog() *repoBacklog {
	return &repoBacklog{}
}

func (b *repoBacklog) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&b.ID, &b.Description, &b.Done)
}

func (b repoBacklog) GetID() int {
	return b.ID
}

const backlogTable = "Задолженность"
const backlogFields = "id,описание,ликвидирована"
const backlogStudentFK = "id_студента"
const backlogDisciplineFK = "id_дисциплины"

func (b repoBacklog) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  backlogTable,
		Fields: backlogFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: backlogStudentFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoStudent() },
			},
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: backlogDisciplineFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoDiscipline() },
			},
		},
	}
}

func (c *repoBacklog) toModel() *models.Backlog {
	if c.model == nil {
		c.model = &models.Backlog{
			ID:          c.ID,
			Description: c.Description,
			Done:        c.Done,
		}

		c.model.Discipline = c.Discipline.toModel()
		c.model.Student = c.Student.toModel()
	}

	return c.model
}

func (s *repoBacklog) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoDiscipline:
		s.Discipline = dep
	case *repoStudent:
		s.Student = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetBacklog(ctx context.Context, id int) (*models.Backlog, error) {
	backlog := NewRepoBacklog()
	filler, err := pgorm.MakeFiller(ctx, r.db, backlogFields, backlogTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrBacklogNotFound
	}

	err = filler.Fill(backlog)

	return backlog.toModel(), nil
}

func (r *DBAisRepository) GetAllBacklogs(ctx context.Context) ([]*models.Backlog, error) {
	backlogs := make([]*models.Backlog, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, backlogFields, backlogTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoBacklog := NewRepoBacklog()
		err = filler.Fill(newRepoBacklog)
		if err != nil {
			return nil, err
		}
		backlogs = append(backlogs, newRepoBacklog.toModel())
	}

	return backlogs, nil
}

// const insertQuery

func (r *DBAisRepository) CreateBacklog(ctx context.Context, description string, disciplineID, studentID int) (*models.Backlog, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO Задолженность(описание,id_дисциплины,id_студента) VALUES ( $1, $2, $3 ) RETURNING id`,
		description, disciplineID, studentID,
	)

	var newID int
	err := row.Scan(&newID)

	if err != nil {
		return nil, err
	}

	return r.GetBacklog(ctx, newID)
}
