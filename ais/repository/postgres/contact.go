package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Контакт(
//     id SERIAL,
//     id_типа_контакта int NOT NULL references ТипКонтакта(id) ON DELETE CASCADE,
//     id_студента int NOT NULL references Студент(id) ON DELETE CASCADE,
//     значение varchar(100) NOT NULL,
//     CONSTRAINT контакт_pk PRIMARY KEY (id)
// );
type Contact struct {
	ID  int
	Def string

	*ContactType
	*Student
}

const contactTable = "Контакт"
const contactIDField = "id"
const contactFields = "id,значение"
const contactStudentFK = "id_студента"
const contactContactTypeFK = "id_типа_контакта"

func (c *Contact) toModel(studentRef *models.Student) *models.Contact {
	contact := &models.Contact{
		ID:          c.ID,
		Def:         c.Def,
		ContactType: c.ContactType.toModel(),
		Student:     studentRef,
	}

	if contact.Student == nil {
		contact.Student = c.Student.toModel(contact)
	}

	return contact
}

func NewPostgresContact(scannable postgres.Scannable) (*Contact, error) {
	contact := &Contact{}

	err := scannable.Scan(&contact.ID, &contact.Def)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrContactNotFound
		}
		return nil, err
	}

	return contact, nil
}

func (c *Contact) Associate(ctx context.Context, r DBAisRepository, studentRef *Student) error {
	contactTypeRow := r.db.QueryRowContext(
		ctx,
		postgres.MakeJoinQuery(
			contactTypeTable, contactTypeFields, contactTypeIDField,
			contactTable, contactContactTypeFK, contactIDField),
		c.ID,
	)

	contactType, err := NewPostgresContactType(contactTypeRow)

	if err != nil {
		return err
	}

	c.ContactType = contactType

	studentRow := r.db.QueryRowContext(
		ctx,
		postgres.MakeJoinQuery(studentTable, studentFields, "id", contactTable, "id_студента", "id"),
		c.ID,
	)

	if studentRef == nil {
		student, err := NewPostgresStudent(studentRow)

		if err != nil {
			return err
		}

		student.Associate(ctx, r, c)
		c.Student = student
	} else {
		c.Student = studentRef
	}

	return nil
}

func makeContactModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Contact, error) {
	contact, err := NewPostgresContact(scannable)

	if err != nil {
		return nil, err
	}

	err = contact.Associate(ctx, r, nil)

	if err != nil {
		return nil, err
	}

	return contact.toModel(nil), nil
}

func (r DBAisRepository) GetContact(ctx context.Context, contactID int) (*models.Contact, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", contactFields, contactTable), contactID)

	return makeContactModel(ctx, r, row)
}

func (r DBAisRepository) GetAllContacts(ctx context.Context) ([]*models.Contact, error) {
	errValue := []*models.Contact{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", contactFields, contactTable))

	if err != nil {
		return errValue, err
	}

	contacts := []*models.Contact{}
	for rows.Next() {
		contact, err := makeContactModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}
