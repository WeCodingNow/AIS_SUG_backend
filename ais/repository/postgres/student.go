package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE Студент(
//     id SERIAL,
//     id_группы int NOT NULL references Группа(id) ON DELETE CASCADE,
//     имя varchar(50) NOT NULL,
//     фамилия varchar(50) NOT NULL,
//     отчество varchar(50),
//     CONSTRAINT студент_pk PRIMARY KEY (id)
// );

type Student struct {
	ID         int
	Name       string
	SecondName string
	ThirdName  sql.NullString

	Contacts []*Contact
	// Group *Group
	// Residence *Residence
}

const studentTable = "Студент"
const studentIDField = "id"
const studentFields = "id,имя,фамилия,отчество"
const studentGroupFK = "id_группы"
const studentResidenceFK = "id_места_жительства"

func (s *Student) toModel(contactRef *models.Contact) *models.Student {
	student := &models.Student{
		ID:         s.ID,
		Name:       s.Name,
		SecondName: s.SecondName,
	}

	if s.ThirdName.Valid {
		student.ThirdName = new(string)
		*student.ThirdName = s.ThirdName.String
	}

	contacts := make([]*models.Contact, 0)

	for _, contact := range s.Contacts {
		if contactRef != nil {
			if contact.ID == contactRef.ID {
				contacts = append(contacts, contactRef)
			} else {
				contacts = append(contacts, contact.toModel(student))
			}
		} else {
			contacts = append(contacts, contact.toModel(student))
		}
	}

	student.Contacts = contacts

	// Group: s.Group.toModel(),
	// Residence: s.Residence.toModel(),
	// Contacts:   contacts,

	return student
}

func NewPostgresStudent(scannable postgres.Scannable) (*Student, error) {
	student := &Student{}

	err := scannable.Scan(&student.ID, &student.Name, &student.SecondName, &student.ThirdName)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrStudentNotFound
		}
		return nil, err
	}

	return student, nil
}

func (s *Student) Associate(ctx context.Context, r DBAisRepository, contactRef *Contact) error {
	contactsRow, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(contactTable, contactFields, contactStudentFK, studentTable, studentIDField, studentIDField),
		s.ID,
	)

	if err != nil {
		return err
	}

	contacts := make([]*Contact, 0)
	for contactsRow.Next() {
		contact, err := NewPostgresContact(contactsRow)

		if err != nil {
			return err
		}

		if contactRef == nil {
			contact.Associate(ctx, r, s)
		} else {
			if contactRef.ID == contact.ID {
				contact = contactRef
			} else {
				contact.Associate(ctx, r, s)
			}
		}

		contacts = append(contacts, contact)
	}

	s.Contacts = contacts

	return nil
}

func makeStudentModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Student, error) {
	student, err := NewPostgresStudent(scannable)

	if err != nil {
		return nil, err
	}

	err = student.Associate(ctx, r, nil)

	if err != nil {
		return nil, err
	}

	return student.toModel(nil), nil
}
func (r DBAisRepository) GetStudent(ctx context.Context, studentID int) (*models.Student, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", studentFields, studentTable), studentID)
	return makeStudentModel(ctx, r, row)
}

func (r DBAisRepository) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	errValue := []*models.Student{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", studentFields, studentTable))

	if err != nil {
		return errValue, err
	}

	students := []*models.Student{}
	for rows.Next() {
		student, err := makeStudentModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		students = append(students, student)
	}

	return students, nil
}
