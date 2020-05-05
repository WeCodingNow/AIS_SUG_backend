package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
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

	model *models.Discipline
}

func NewRepoDiscipline() *repoDiscipline {
	return &repoDiscipline{
		ControlEvents: make(map[int]*repoControlEvent),
	}
}

func (s *repoDiscipline) Fill(scannable Scannable) {
	scannable.Scan(&s.ID, &s.Name, &s.Hours)
}

func (s repoDiscipline) GetID() int {
	return s.ID
}

const disciplineFields = "id,название,часы"
const disciplineTable = "Дисциплина"

func (c repoDiscipline) GetDescription() ModelDescription {
	return ModelDescription{
		Table:  disciplineTable,
		Fields: disciplineFields,
		Dependencies: []ModelDependency{
			{
				DependencyType:     OneToMany,
				DepForeignKeyField: controlEventDisciplineFK,
				ModelMaker:         func() RepoModel { return NewRepoControlEvent() },
			},
		},
	}
}

func (c *repoDiscipline) toModel() *models.Discipline {
	if c.model == nil {
		c.model = &models.Discipline{
			ID:    c.ID,
			Name:  c.Name,
			Hours: c.Hours,
		}

		controlEvents := make([]*models.ControlEvent, 0, len(c.ControlEvents))
		for _, repoM := range c.ControlEvents {
			controlEvents = append(controlEvents, repoM.toModel())
		}
		c.model.ControlEvents = controlEvents
	}

	return c.model
}

func (s *repoDiscipline) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoControlEvent:
		s.ControlEvents[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetDiscipline(ctx context.Context, id int) (*models.Discipline, error) {
	discipline := NewRepoDiscipline()
	filler, err := MakeFiller(ctx, r.db, disciplineFields, disciplineTable, &id)

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
	filler, err := MakeFiller(ctx, r.db, disciplineFields, disciplineTable, nil)

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
