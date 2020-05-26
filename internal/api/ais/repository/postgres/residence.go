package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/utils"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE МестоЖительства(
//     id SERIAL,
//     адрес varchar(100) NOT NULL,
//     город varchar(20) NOT NULL,
//     общежитие boolean NOT NULL,
//     CONSTRAINT место_жительства_pk PRIMARY KEY (id)
// );

type repoResidence struct {
	ID        int
	Address   string
	City      string
	Community bool

	Students map[int]*repoStudent

	model *models.Residence
}

func NewRepoResidence() *repoResidence {
	return &repoResidence{
		Students: make(map[int]*repoStudent),
	}
}

func (s *repoResidence) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Address, &s.City, &s.Community)
}

func (s repoResidence) GetID() int {
	return s.ID
}

const residenceTable = "МестоЖительства"
const residenceIDField = "id"
const residenceFields = "id,адрес,город,общежитие"

func (c repoResidence) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  residenceTable,
		Fields: residenceFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: studentResidenceFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoStudent() },
			},
		},
	}
}

func (c *repoResidence) toModel() *models.Residence {
	if c.model == nil {
		c.model = &models.Residence{
			ID:        c.ID,
			Address:   c.Address,
			City:      c.City,
			Community: c.Community,
		}

		students := make([]*models.Student, 0, len(c.Students))
		for _, repoS := range c.Students {
			students = append(students, repoS.toModel())
		}
		c.model.Students = students
	}

	return c.model
}

func (s *repoResidence) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoStudent:
		s.Students[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetResidence(ctx context.Context, id int) (*models.Residence, error) {
	residence := NewRepoResidence()
	filler, err := pgorm.MakeFiller(ctx, r.db, residenceFields, residenceTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrResidenceNotFound
	}

	err = filler.Fill(residence)

	return residence.toModel(), nil
}

func (r *DBAisRepository) GetAllResidences(ctx context.Context) ([]*models.Residence, error) {
	residences := make([]*models.Residence, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, residenceFields, residenceTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoResidence := NewRepoResidence()
		err = filler.Fill(newRepoResidence)
		if err != nil {
			return nil, err
		}
		residences = append(residences, newRepoResidence.toModel())
	}

	return residences, nil
}

func (r *DBAisRepository) CreateResidence(ctx context.Context, address, city string, community bool) (*models.Residence, error) {
	residenceRow := r.db.QueryRowContext(ctx, utils.MakeInsertQueryReturningModel(
		residenceTable, strings.Split(residenceFields, ",")[1:]...),
		address, city, community,
	)

	repoResidence := new(repoResidence)

	err := repoResidence.Fill(residenceRow)

	if err != nil {
		return nil, err
	}

	return repoResidence.toModel(), nil
}
