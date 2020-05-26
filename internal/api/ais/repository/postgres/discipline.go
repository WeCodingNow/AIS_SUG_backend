package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE Дисциплина(
//     id SERIAL,
//     название varchar(150),
//     часы int,
//     CONSTRAINT дисциплина_pk PRIMARY KEY (id)
// );

type repoDiscipline struct {
	ID    int
	Name  string
	Hours int

	ControlEvents map[int]*repoControlEvent
	Semesters     map[int]*repoSemester
	Backlogs      map[int]*repoBacklog

	model *models.Discipline
}

func NewRepoDiscipline() *repoDiscipline {
	return &repoDiscipline{
		ControlEvents: make(map[int]*repoControlEvent),
		Semesters:     make(map[int]*repoSemester),
		Backlogs:      make(map[int]*repoBacklog),
	}
}

func (s *repoDiscipline) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Name, &s.Hours)
}

func (s repoDiscipline) GetID() int {
	return s.ID
}

const disciplineFields = "id,название,часы"
const disciplineTable = "Дисциплина"

func (c repoDiscipline) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  disciplineTable,
		Fields: disciplineFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: controlEventDisciplineFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoControlEvent() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: backlogDisciplineFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoBacklog() },
			},
			{
				DependencyType:     pgorm.ManyToMany,
				ForeignKeyField:    disciplineSemesterMTMDisciplineKey,
				DepForeignKeyField: disciplineSemesterMTMSemesterKey,
				ManyToManyTable:    disciplineSemesterMtM,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoSemester() },
			},
		},
	}
}

func (d *repoDiscipline) toModel() *models.Discipline {
	if d.model == nil {
		d.model = &models.Discipline{
			ID:    d.ID,
			Name:  d.Name,
			Hours: d.Hours,
		}

		controlEvents := make([]*models.ControlEvent, 0, len(d.ControlEvents))
		for _, repoM := range d.ControlEvents {
			controlEvents = append(controlEvents, repoM.toModel())
		}
		d.model.ControlEvents = controlEvents

		semesters := make([]*models.Semester, 0, len(d.Semesters))
		for _, repoS := range d.Semesters {
			semesters = append(semesters, repoS.toModel())
		}
		d.model.Semesters = semesters

		backlogs := make([]*models.Backlog, 0, len(d.Backlogs))
		for _, repoB := range d.Backlogs {
			backlogs = append(backlogs, repoB.toModel())
		}
		d.model.Backlogs = backlogs
	}

	return d.model
}

func (d *repoDiscipline) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoControlEvent:
		d.ControlEvents[dep.ID] = dep
	case *repoSemester:
		d.Semesters[dep.ID] = dep
	case *repoBacklog:
		d.Backlogs[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetDiscipline(ctx context.Context, id int) (*models.Discipline, error) {
	discipline := NewRepoDiscipline()
	filler, err := pgorm.MakeFiller(ctx, r.db, disciplineFields, disciplineTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrDisciplineNotFound
	}

	err = filler.Fill(discipline)

	return discipline.toModel(), nil
}

func (r *DBAisRepository) GetAllDisciplines(ctx context.Context) ([]*models.Discipline, error) {
	disciplines := make([]*models.Discipline, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, disciplineFields, disciplineTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoDiscipline := NewRepoDiscipline()
		err = filler.Fill(newRepoDiscipline)
		if err != nil {
			return nil, err
		}
		disciplines = append(disciplines, newRepoDiscipline.toModel())
	}

	return disciplines, nil
}
