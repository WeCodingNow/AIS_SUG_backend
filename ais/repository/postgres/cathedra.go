package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Кафедра(
//     id SERIAL,
//     название varchar(100) NOT NULL UNIQUE,
//     короткое_название varchar(10) NOT NULL UNIQUE,
//     CONSTRAINT кафедра_pk PRIMARY KEY (id)
// );

type Cathedra struct {
	ID        int
	Name      string
	ShortName string

	Groups []*Group
}

const cathedraTable = "Кафедра"
const cathedraIDField = "id"
const cathedraFields = "id,название,короткое_название"

func (c *Cathedra) toModel(groupRef *models.Group) *models.Cathedra {
	cathedra := &models.Cathedra{
		ID:        c.ID,
		Name:      c.Name,
		ShortName: c.ShortName,
	}

	groups := make([]*models.Group, 0)

	for _, group := range c.Groups {
		if groupRef != nil {
			if group.ID == groupRef.ID {
				groups = append(groups, groupRef)
			} else {
				groups = append(groups, group.toModel(nil, cathedra, nil))
			}
		} else {
			groups = append(groups, group.toModel(nil, cathedra, nil))
		}
	}

	cathedra.Groups = groups

	return cathedra
}

func NewPostgresCathedra(scannable postgres.Scannable) (*Cathedra, error) {
	cathedra := &Cathedra{}

	err := scannable.Scan(&cathedra.ID, &cathedra.Name, &cathedra.ShortName)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrCathedraNotFound
		}
		return nil, err
	}

	return cathedra, nil
}

func (c *Cathedra) Associate(ctx context.Context, r DBAisRepository, groupRef *Group) error {
	groupRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(groupTable, groupFields, groupCathedraFK, cathedraTable, cathedraIDField, cathedraIDField),
		c.ID,
	)

	if err != nil {
		return err
	}

	groups := make([]*Group, 0)
	for groupRows.Next() {
		group, err := NewPostgresGroup(groupRows)

		if err != nil {
			return err
		}

		if groupRef == nil {
			group.Associate(ctx, r, nil, c, nil)
		} else {
			if groupRef.ID == group.ID {
				group = groupRef
			} else {
				group.Associate(ctx, r, nil, c, nil)
			}
		}

		groups = append(groups, group)
	}

	c.Groups = groups

	return nil
}

func makeCathedraModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Cathedra, error) {
	cathedra, err := NewPostgresCathedra(scannable)

	if err != nil {
		return nil, err
	}

	err = cathedra.Associate(ctx, r, nil)

	if err != nil {
		return nil, err
	}

	return cathedra.toModel(nil), nil
}

func (r DBAisRepository) GetCathedra(ctx context.Context, cathedraID int) (*models.Cathedra, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", cathedraFields, cathedraTable), cathedraID)
	return makeCathedraModel(ctx, r, row)
}

func (r DBAisRepository) GetAllCathedras(ctx context.Context) ([]*models.Cathedra, error) {
	errValue := []*models.Cathedra{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", cathedraFields, cathedraTable))

	if err != nil {
		return errValue, err
	}

	cathedras := []*models.Cathedra{}
	for rows.Next() {
		cathedra, err := makeCathedraModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		cathedras = append(cathedras, cathedra)
	}

	return cathedras, nil
}
