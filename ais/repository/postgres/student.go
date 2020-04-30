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

	Contacts  []*Contact
	Group     *Group
	Residence *Residence
	Marks     []*Mark
}

const studentTable = "Студент"
const studentIDField = "id"
const studentFields = "id,имя,фамилия,отчество"
const studentGroupFK = "id_группы"
const studentResidenceFK = "id_места_жительства"

func (s *Student) toModel(
	contactRef *models.Contact, groupRef *models.Group,
	residenceRef *models.Residence, markRef *models.Mark,
) *models.Student {
	student := &models.Student{
		ID:         s.ID,
		Name:       s.Name,
		SecondName: s.SecondName,

		Group:     groupRef,
		Residence: residenceRef,
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

	if student.Group == nil {
		student.Group = s.Group.toModel(student, nil, nil)
	}

	if student.Residence == nil {
		student.Residence = s.Residence.toModel(student)
	}

	marks := make([]*models.Mark, 0)
	for _, mark := range s.Marks {
		if markRef != nil {
			if mark.ID == markRef.ID {
				marks = append(marks, markRef)
			} else {
				marks = append(marks, mark.toModel(markRef.ControlEvent, student))
			}
		} else {
			marks = append(marks, mark.toModel(nil, student))
		}
	}
	student.Marks = marks

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

func (s *Student) Associate(
	ctx context.Context, r DBAisRepository,
	contactRef *Contact, groupRef *Group, residenceRef *Residence, markRef *Mark,
) error {
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

	if groupRef == nil {
		groupRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(groupTable, groupFields, "id", studentTable, studentGroupFK, "id"),
			s.ID,
		)

		group, err := NewPostgresGroup(groupRow)

		if err != nil {
			return err
		}

		group.Associate(ctx, r, s, nil, nil)
		s.Group = group
	} else {
		s.Group = groupRef
	}

	if residenceRef == nil {
		residenceRow := r.db.QueryRowContext(
			ctx,
			postgres.MakeJoinQuery(residenceTable, residenceFields, "id", studentTable, studentResidenceFK, "id"),
			s.ID,
		)

		residence, err := NewPostgresResidence(residenceRow)

		if err != nil {
			return err
		}

		residence.Associate(ctx, r, s)
		s.Residence = residence
	} else {
		s.Residence = residenceRef
	}

	markRows, err := r.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(markTable, markFields, markStudentFK, studentTable, studentIDField, studentIDField),
		s.ID,
	)

	if err != nil {
		return err
	}

	marks := make([]*Mark, 0)
	for markRows.Next() {
		mark, err := NewPostgresMark(markRows)

		if err != nil {
			return err
		}

		if markRef == nil {
			mark.Associate(ctx, r, nil, s)
			// mark.Student = s
		} else {
			if markRef.ID == mark.ID {
				mark = markRef
			} else {
				mark.Associate(ctx, r, nil, s)
			}
		}
		mark.Student = s
		marks = append(marks, mark)
	}

	s.Marks = marks

	return nil
}

func makeStudentModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Student, error) {
	student, err := NewPostgresStudent(scannable)

	if err != nil {
		return nil, err
	}

	err = student.Associate(ctx, r, nil, nil, nil, nil)

	if err != nil {
		return nil, err
	}

	return student.toModel(nil, nil, nil, nil), nil
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
