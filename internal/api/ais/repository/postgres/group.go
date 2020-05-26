package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

type repoGroup struct {
	ID     int
	Number int

	Cathedra  *repoCathedra
	Students  map[int]*repoStudent
	Semesters map[int]*repoSemester
	model     *models.Group
}

func NewRepoGroup() *repoGroup {
	return &repoGroup{
		Students:  make(map[int]*repoStudent),
		Semesters: make(map[int]*repoSemester),
	}
}

func (g *repoGroup) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&g.ID, &g.Number)
}

func (g repoGroup) GetID() int {
	return g.ID
}

const groupTable = "Группа"
const groupIDField = "id"
const groupFields = "id,номер"

const groupCathedraFK = "id_кафедры"
const groupSemesterMtM = "Группа__Семестр"
const groupSemesterMTMGroupKey = "id_группы"
const groupSemesterMTMSemesterKey = "id_семестра"

func (c repoGroup) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  groupTable,
		Fields: groupFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:  pgorm.ManyToOne,
				ForeignKeyField: groupCathedraFK,
				ModelMaker:      func() pgorm.RepoModel { return NewRepoCathedra() },
			},
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: studentGroupFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoStudent() },
			},
			{
				DependencyType:     pgorm.ManyToMany,
				ForeignKeyField:    groupSemesterMTMGroupKey,
				DepForeignKeyField: groupSemesterMTMSemesterKey,
				ManyToManyTable:    groupSemesterMtM,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoSemester() },
			},
		},
	}
}

func (c *repoGroup) toModel() *models.Group {
	if c.model == nil {
		c.model = &models.Group{
			ID:       c.ID,
			Number:   c.Number,
			Cathedra: c.Cathedra.toModel(),
		}

		c.model.Cathedra = c.Cathedra.toModel()

		students := make([]*models.Student, 0, len(c.Students))
		for _, repoS := range c.Students {
			students = append(students, repoS.toModel())
		}
		c.model.Students = students

		semesters := make([]*models.Semester, 0, len(c.Semesters))
		for _, repoS := range c.Semesters {
			semesters = append(semesters, repoS.toModel())
		}
		c.model.Semesters = semesters
	}

	return c.model
}

func (c *repoGroup) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoCathedra:
		c.Cathedra = dep
	case *repoStudent:
		c.Students[dep.ID] = dep
	case *repoSemester:
		c.Semesters[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	group := NewRepoGroup()
	filler, err := pgorm.MakeFiller(ctx, r.db, groupFields, groupTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrGroupNotFound
	}

	err = filler.Fill(group)

	return group.toModel(), nil
}

func (r *DBAisRepository) GetAllGroups(ctx context.Context) ([]*models.Group, error) {
	groups := make([]*models.Group, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, groupFields, groupTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoGroup := NewRepoGroup()
		err = filler.Fill(newRepoGroup)
		if err != nil {
			return nil, err
		}
		groups = append(groups, newRepoGroup.toModel())
	}

	return groups, nil
}
