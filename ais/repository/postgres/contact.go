package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
)

type Contact struct {
	ID            int
	ContactTypeID int
	StudentID     int
	// *ContactType
	// *Student
	Def string
}

// CREATE TABLE Кафедра(
//     id SERIAL,
//     название varchar(100) NOT NULL UNIQUE,
//     короткое_название varchar(10) NOT NULL UNIQUE,
//     CONSTRAINT кафедра_pk PRIMARY KEY (id)
// );

func toPostgresContact(c *models.Contact) *Contact {
	// contactType, err := r.GetContactType(ctx, c.ContactType.ID)
	// c.ContactType.ID
	return &Contact{
		ID:            c.ID,
		ContactTypeID: c.ContactType.ID,
		StudentID:     c.Student.ID,
		Def:           c.Def,
		// toPostgresContactType(),
		// c.
		// c.Def,
	}
}

func toModelContact(r DBAisRepository, ctx context.Context, c *Contact) *models.Contact {
	contactType, err := r.GetContactType(ctx, c.ContactTypeID)

	if err != nil {
		panic(err)
	}

	// student, err := r.GetStudent(ctx, c.StudentID)

	// if err != nil {
	// 	panic(err)
	// }

	return &models.Contact{
		ID:          c.ID,
		ContactType: contactType,
		// Student: student,
		Def: c.Def,
	}
}

// const createCathedraQuery = `INSERT INTO Кафедра(название, короткое_название) VALUES ( $1, $2 )`

// func (r AisRepository) CreateCathedra(ctx context.Context, name, shortName string) error {
// 	_, err := r.db.ExecContext(ctx, createCathedraQuery,
// 		name, shortName,
// 	)

// 	return err
// }

const getContactQuery = `SELECT id, id_типа_контакта, id_студента, значение FROM Контакт WHERE id = $1`

func (r DBAisRepository) GetContact(ctx context.Context, contactID int) (*models.Contact, error) {
	row := r.db.QueryRowContext(ctx, getContactQuery, contactID)

	contact := new(Contact)
	err := row.Scan(&contact.ID, &contact.ContactTypeID, &contact.StudentID, &contact.Def)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrContactNotFound
		}
		return nil, err
	}

	return toModelContact(r, ctx, contact), nil
}

const getAllContactsQuery = `SELECT * FROM Контакт`

func (r DBAisRepository) GetAllContacts(ctx context.Context) ([]*models.Contact, error) {
	rows, err := r.db.QueryContext(ctx, getAllContactsQuery)
	contacts := make([]*models.Contact, 0)

	if err != nil {
		return contacts, err
	}

	for rows.Next() {
		contact := new(Contact)
		if err := rows.Scan(&contact.ID, &contact.ContactTypeID, &contact.StudentID, &contact.Def); err != nil {
			return []*models.Contact{}, err
		}
		contacts = append(contacts, toModelContact(r, ctx, contact))
	}

	return contacts, nil
}
