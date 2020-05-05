package postgres

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type repoContactType struct {
	ID  int
	Def string

	model *models.ContactType
}

func NewRepoContactType() *repoContactType {
	return &repoContactType{}
}

func (s *repoContactType) Fill(scannable Scannable) {
	scannable.Scan(&s.ID, &s.Def)
}

func (s repoContactType) GetID() int {
	return s.ID
}

const contactTypeFields = "id,обозначение"
const contactTypeTable = "ТипКонтакта"

func (c repoContactType) GetDescription() ModelDescription {
	return ModelDescription{
		Table:        contactTypeTable,
		Fields:       contactTypeFields,
		Dependencies: []ModelDependency{},
	}
}

func (c *repoContactType) toModel() *models.ContactType {
	if c.model == nil {
		c.model = &models.ContactType{
			ID:  c.ID,
			Def: c.Def,
		}
	}

	return c.model
}

func (s *repoContactType) AcceptDep(dep interface{}) error {
	return nil
}

func (r *DBAisRepository) GetContactType(ctx context.Context, id int) (*models.ContactType, error) {
	contactType := NewRepoContactType()
	filler, err := MakeFiller(ctx, r.db, contactTypeFields, contactTypeTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrContactTypeNotFound
	}

	err = filler.Fill(contactType)

	return contactType.toModel(), nil
}

func (r *DBAisRepository) GetAllContactTypes(ctx context.Context) ([]*models.ContactType, error) {
	contactTypes := make([]*models.ContactType, 0)
	filler, err := MakeFiller(ctx, r.db, contactTypeFields, contactTypeTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepocontactType := NewRepoContactType()
		err = filler.Fill(newRepocontactType)
		if err != nil {
			return nil, err
		}
		contactTypes = append(contactTypes, newRepocontactType.toModel())
	}

	return contactTypes, nil
}
