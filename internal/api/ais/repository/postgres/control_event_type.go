package postgres

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

// CREATE TABLE ТипКонтрольногоМероприятия(
//     id SERIAL,
//     обозначение varchar(50),
//     CONSTRAINT тип_контрольного_мероприятие_pk PRIMARY KEY (id)
// );

type repoControlEventType struct {
	ID  int
	Def string

	model *models.ControlEventType
}

func NewRepoControlEventType() *repoControlEventType {
	return &repoControlEventType{}
}

func (s *repoControlEventType) Fill(scannable pgorm.Scannable) error {
	return scannable.Scan(&s.ID, &s.Def)
}

func (s repoControlEventType) GetID() int {
	return s.ID
}

const controlEventTypeFields = "id,обозначение"
const controlEventTypeTable = "ТипКонтрольногоМероприятия"

func (c repoControlEventType) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:        controlEventTypeTable,
		Fields:       controlEventTypeFields,
		Dependencies: []pgorm.ModelDependency{},
	}
}

func (c *repoControlEventType) toModel() *models.ControlEventType {
	if c.model == nil {
		c.model = &models.ControlEventType{
			ID:  c.ID,
			Def: c.Def,
		}
	}

	return c.model
}

func (s *repoControlEventType) AcceptDep(dep interface{}) error {
	return nil
}

func (r *DBAisRepository) GetControlEventType(ctx context.Context, id int) (*models.ControlEventType, error) {
	controlEventType := NewRepoControlEventType()
	filler, err := pgorm.MakeFiller(ctx, r.db, controlEventTypeFields, controlEventTypeTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrControlEventTypeNotFound
	}

	err = filler.Fill(controlEventType)

	return controlEventType.toModel(), nil
}

func (r *DBAisRepository) GetAllControlEventTypes(ctx context.Context) ([]*models.ControlEventType, error) {
	controlEventTypes := make([]*models.ControlEventType, 0)
	filler, err := pgorm.MakeFiller(ctx, r.db, controlEventTypeFields, controlEventTypeTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoControlEventType := NewRepoControlEventType()
		err = filler.Fill(newRepoControlEventType)
		if err != nil {
			return nil, err
		}
		controlEventTypes = append(controlEventTypes, newRepoControlEventType.toModel())
	}

	return controlEventTypes, nil
}
