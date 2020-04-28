package postgres

import (
	"context"
	"database/sql"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
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
	ID          int
	GroupID     int
	ResidenceID int
	ContactIDs  []int
	Name        string
	SecondName  string
	ThirdName   sql.NullString
}

func toPostgresStudent(g *models.Student) *Student {
	return nil
}

func toModelStudent(r DBAisRepository, ctx context.Context, g *Student) (*models.Student, error) {
	residenceModel, err := r.GetResidence(ctx, g.ResidenceID)
	if err != nil {
		return nil, err
	}

	contactModels := make([]*models.Contact, 0, len(g.ContactIDs))

	for _, contactID := range g.ContactIDs {
		contact, err := r.GetContact(ctx, contactID)
		if err != nil {
			return nil, err
		}

		contactModels = append(contactModels, contact)
	}

	retStudent := &models.Student{
		ID:         g.ID,
		Group:      nil,
		Residence:  residenceModel,
		Contacts:   contactModels,
		Name:       g.Name,
		SecondName: g.SecondName,
	}

	if g.ThirdName.Valid {
		retStudent.ThirdName = new(string)
		*retStudent.ThirdName = g.ThirdName.String
	}

	retStudent.Group, err = r.GetGroupRecursive(ctx, g.GroupID, retStudent)
	if err != nil {
		return nil, err
	}

	return retStudent, nil
}

const getStudentQuery = `
SELECT id, id_группы, id_места_жительства, имя, фамилия, отчество
	FROM Студент
	WHERE id = $1`

const getContactIDsInStudentQuery = `
SELECT к.id
	FROM Студент as с
	JOIN Контакт as к
	ON с.id = к.id_студента
	WHERE с.id = $1`

func (s *Student) Fill(ctx context.Context, db *sql.DB, sc Scannable) (*Student, error) {
	err := s.hydrate(sc)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ais.ErrStudentNotFound
		}
		return nil, err
	}

	idRows, err := db.QueryContext(ctx, getContactIDsInStudentQuery, s.ID)

	if err != nil {
		return nil, err
	}

	for idRows.Next() {
		var contactID int
		idRows.Scan(&contactID)
		s.ContactIDs = append(s.ContactIDs, contactID)
	}

	return s, nil
}

func (s *Student) hydrate(sc Scannable) error {
	return sc.Scan(&s.ID, &s.GroupID, &s.ResidenceID, &s.Name, &s.SecondName, &s.ThirdName)
}

func (r DBAisRepository) GetStudent(ctx context.Context, studentID int) (*models.Student, error) {
	row := r.db.QueryRowContext(ctx, getStudentQuery, studentID)
	student := new(Student)
	student, err := student.Fill(ctx, r.db, row)

	if err != nil {
		return nil, err
	}

	studentModel, err := toModelStudent(r, ctx, student)
	if err != nil {
		return nil, err
	}

	return studentModel, nil
}

const getAllStudentsQuery = `
SELECT id, id_группы, id_места_жительства, имя, фамилия, отчество
	FROM Студент`

func (r DBAisRepository) GetAllStudents(ctx context.Context) ([]*models.Student, error) {
	rows, err := r.db.QueryContext(ctx, getAllStudentsQuery)
	students := make([]*models.Student, 0)

	if err != nil {
		return students, err
	}

	for rows.Next() {
		student := new(Student)
		student.Fill(ctx, r.db, rows)

		studentModel, err := toModelStudent(r, ctx, student)
		if err != nil {
			return []*models.Student{}, err
		}
		students = append(students, studentModel)
	}

	return students, nil
}
