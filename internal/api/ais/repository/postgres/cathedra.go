package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

type repoCathedra struct {
	ID        int
	Name      string
	ShortName string

	Groups map[int]*repoGroup

	model *models.Cathedra
}

func NewRepoCathedra() *repoCathedra {
	return &repoCathedra{
		Groups: make(map[int]*repoGroup),
	}
}

func (c *repoCathedra) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&c.ID, &c.Name, &c.ShortName)
}

func (c repoCathedra) GetID() int {
	return c.ID
}

const cathedraTable = "Кафедра"
const cathedraFields = "id,название,короткое_название"

func (c repoCathedra) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  cathedraTable,
		Fields: cathedraFields,
		Dependencies: []pgorm.ModelDependency{
			{
				DependencyType:     pgorm.OneToMany,
				DepForeignKeyField: groupCathedraFK,
				ModelMaker:         func() pgorm.RepoModel { return NewRepoGroup() },
			},
		},
	}
}

func (c *repoCathedra) toModel() *models.Cathedra {
	if c.model == nil {
		c.model = &models.Cathedra{
			ID:        c.ID,
			Name:      c.Name,
			ShortName: c.ShortName,
		}

		groups := make([]*models.Group, 0, len(c.Groups))
		for _, repoG := range c.Groups {
			groups = append(groups, repoG.toModel())
		}

		c.model.Groups = groups
	}

	return c.model
}

func (c *repoCathedra) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoGroup:
		c.Groups[dep.ID] = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetCathedra(ctx context.Context, id int) (*models.Cathedra, error) {
	cathedra := NewRepoCathedra()
	filler, err := pgorm.MakeFiller(ctx, r.db, cathedraFields, cathedraTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrCathedraNotFound
	}

	err = filler.Fill(cathedra)

	return cathedra.toModel(), nil
}

func (r *DBAisRepository) GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error) {
	cathedras := make([]*models.Cathedra, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, cathedraFields, cathedraTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoCathedra := NewRepoCathedra()
		err = filler.Fill(newRepoCathedra)
		if err != nil {
			return nil, err
		}
		cathedras = append(cathedras, newRepoCathedra.toModel())
	}

	return cathedras, nil
}
