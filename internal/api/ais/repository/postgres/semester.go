package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE Семестр(
//     id SERIAL,
//     номер int NOT NULL,
//     начало date NOT NULL,
//     конец date,
//     CONSTRAINT семестр_pk PRIMARY KEY (id)
// );

type repoSemester struct {
	ID        int
	Number    int
	Beginning sql.NullTime
	End       sql.NullTime

	Groups        map[int]*repoGroup
	ControlEvents map[int]*repoControlEvent
	Disciplines   map[int]*repoDiscipline
	model         *models.Semester
}

func NewRepoSemester() *repoSemester {
	return &repoSemester{
		Groups:        make(map[int]*repoGroup),
		ControlEvents: make(map[int]*repoControlEvent),
		Disciplines:   make(map[int]*repoDiscipline),
	}
}

func (s *repoSemester) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Number, &s.Beginning, &s.End)
}

func (s repoSemester) GetID() int {
	return s.ID
}

const semesterTable = "Семестр"
const semesterFields = "id,номер,начало,конец"

const disciplineSemesterMtM = "Дисциплина__Семестр"
const disciplineSemesterMTMDisciplineKey = "id_дисциплины"
const disciplineSemesterMTMSemesterKey = "id_семестра"

func (c repoSemester) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  semesterTable,
		Fields: semesterFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:     pgorm.ManyToMany,
				ForeignKeyField:    groupSemesterMTMSemesterKey,
				DepForeignKeyField: groupSemesterMTMGroupKey,
				ManyToManyTable:    groupSemesterMtM,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoGroup() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: controlEventSemesterFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoControlEvent() },
			},
			{
				DependencyType:     pgorm.ManyToMany,
				ForeignKeyField:    disciplineSemesterMTMSemesterKey,
				DepForeignKeyField: disciplineSemesterMTMDisciplineKey,
				ManyToManyTable:    disciplineSemesterMtM,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoDiscipline() },
			},
		},
	}
}
func (s *repoSemester) toModel() *models.Semester {
	if s.model == nil {
		s.model = &models.Semester{
			ID:        s.ID,
			Number:    s.Number,
			Beginning: s.Beginning.Time,
		}

		if s.End.Valid {
			*s.model.End = s.End.Time
		}

		groups := make([]*models.Group, 0, len(s.Groups))
		for _, repoG := range s.Groups {
			groups = append(groups, repoG.toModel())
		}
		s.model.Groups = groups

		controlEvents := make([]*models.ControlEvent, 0, len(s.Groups))
		for _, repoCe := range s.ControlEvents {
			controlEvents = append(controlEvents, repoCe.toModel())
		}
		s.model.ControlEvents = controlEvents

		disciplines := make([]*models.Discipline, 0, len(s.Disciplines))
		for _, repoD := range s.Disciplines {
			disciplines = append(disciplines, repoD.toModel())
		}
		s.model.Disciplines = disciplines
	}

	return s.model
}

func (s *repoSemester) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoGroup:
		s.Groups[dep.ID] = dep
	case *repoControlEvent:
		s.ControlEvents[dep.ID] = dep
	case *repoDiscipline:
		s.Disciplines[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetSemester(ctx context.Context, id int) (*models.Semester, error) {
	semester := NewRepoSemester()
	filler, err := pgorm.MakeFiller(ctx, r.db, semesterFields, semesterTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrSemesterNotFound
	}

	err = filler.Fill(semester)

	return semester.toModel(), nil
}

func (r *DBAisRepository) GetAllSemesters(ctx context.Context) ([]*models.Semester, error) {
	semesters := make([]*models.Semester, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, semesterFields, semesterTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoSemester := NewRepoSemester()
		err = filler.Fill(newRepoSemester)
		if err != nil {
			return nil, err
		}
		semesters = append(semesters, newRepoSemester.toModel())
	}

	return semesters, nil
}
