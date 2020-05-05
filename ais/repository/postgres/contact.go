package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

// CREATE TABLE Контакт(
//     id SERIAL,
//     id_типа_контакта int NOT NULL references ТипКонтакта(id) ON DELETE CASCADE,
//     id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     значение varchar(100) NOT NULL,
//     CONSTRAINT контакт_pk PRIMARY KEY (id)
// );

type repoContact struct {
	ID  int
	Def string

	ContactType *repoContactType
	Student     *repoStudent

	model *models.Contact
}

func NewRepoContact() *repoContact {
	return &repoContact{}
}

func (s *repoContact) Fill(scannable Scannable) {
	scannable.Scan(&s.ID, &s.Def)
}

func (s repoContact) GetID() int {
	return s.ID
}

const contactTable = "Контакт"
const contactFields = "id,значение"
const contactStudentFK = "id_студента"
const contactContactTypeFK = "id_типа_контакта"

func (c repoContact) GetDescription() ModelDescription {
	return ModelDescription{
		Table:  contactTable,
		Fields: contactFields,
		Dependencies: []ModelDependency{
			{
				DependencyType:  ManyToOne,
				ForeignKeyField: contactStudentFK,
				ModelMaker:      func() RepoModel { return NewRepoStudent() },
			},
			{
				DependencyType:  ManyToOne,
				ForeignKeyField: contactContactTypeFK,
				ModelMaker:      func() RepoModel { return NewRepoContactType() },
			},
		},
	}
}

func (c *repoContact) toModel() *models.Contact {
	if c.model == nil {
		c.model = &models.Contact{
			ID:  c.ID,
			Def: c.Def,
		}

		c.model.ContactType = c.ContactType.toModel()
		c.model.Student = c.Student.toModel()
	}

	return c.model
}

func (s *repoContact) AcceptDep(dep interface{}) error {
	switch dep := dep.(type) {
	case *repoStudent:
		s.Student = dep
	case *repoContactType:
		s.ContactType = dep
	default:
		return fmt.Errorf("no dependency for %v", dep)
	}
	return nil
}

func (r *DBAisRepository) GetContact(ctx context.Context, id int) (*models.Contact, error) {
	contact := NewRepoContact()
	filler, err := MakeFiller(ctx, r.db, contactFields, contactTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, ais.ErrContactNotFound
	}

	err = filler.Fill(contact)

	return contact.toModel(), nil
}

func (r *DBAisRepository) GetAllContacts(ctx context.Context) ([]*models.Contact, error) {
	contacts := make([]*models.Contact, 0)
	filler, err := MakeFiller(ctx, r.db, contactFields, contactTable, nil)

	if err != nil {
		return nil, err
	}

	for filler.Next() {
		newRepoContact := NewRepoContact()
		err = filler.Fill(newRepoContact)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, newRepoContact.toModel())
	}

	return contacts, nil
}
