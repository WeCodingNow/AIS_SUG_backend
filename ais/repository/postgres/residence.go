package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/WeCodingNow/AIS_SUG_backend/utils/delivery/postgres"
)

// CREATE TABLE МестоЖительства(
//     id SERIAL,
//     адрес varchar(100) NOT NULL,
//     город varchar(20) NOT NULL,
//     общежитие boolean NOT NULL,
//     CONSTRAINT место_жительства_pk PRIMARY KEY (id)
// );
type Residence struct {
	ID        int
	Address   string
	City      string
	Community bool
	Students  []*Student
}

const residenceTable = "МестоЖительства"
const residenceIDField = "id"
const residenceFields = "id,адрес,город,общежитие"

func (r *Residence) toModel(studentRef *models.Student) *models.Residence {
	residence := &models.Residence{
		ID:        r.ID,
		Address:   r.Address,
		City:      r.City,
		Community: r.Community,
	}

	students := make([]*models.Student, 0)

	for _, student := range r.Students {
		if studentRef != nil {
			if student.ID == studentRef.ID {
				students = append(students, studentRef)
			} else {
				students = append(students, student.toModel(nil, nil, residence, nil))
			}
		} else {
			students = append(students, student.toModel(nil, nil, residence, nil))
		}
	}

	residence.Students = students

	return residence
}

func NewPostgresResidence(scannable postgres.Scannable) (*Residence, error) {
	residence := &Residence{}

	err := scannable.Scan(&residence.ID, &residence.Address, &residence.City, &residence.Community)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ais.ErrResidenceNotFound
		}
		return nil, err
	}

	return residence, nil
}

func (r *Residence) Associate(ctx context.Context, repo DBAisRepository, studentRef *Student) error {
	studentRows, err := repo.db.QueryContext(
		ctx,
		postgres.MakeJoinQuery(studentTable, studentFields, studentResidenceFK, residenceTable, residenceIDField, residenceIDField),
		r.ID,
	)

	if err != nil {
		return err
	}

	students := make([]*Student, 0)
	for studentRows.Next() {
		student, err := NewPostgresStudent(studentRows)

		if err != nil {
			return err
		}

		if studentRef == nil {
			student.Associate(ctx, repo, nil, nil, r, nil)
		} else {
			if studentRef.ID == student.ID {
				student = studentRef
			} else {
				student.Associate(ctx, repo, nil, nil, r, nil)
			}
		}

		students = append(students, student)
	}

	r.Students = students

	return nil
}

func makeResidenceModel(ctx context.Context, r DBAisRepository, scannable postgres.Scannable) (*models.Residence, error) {
	residence, err := NewPostgresResidence(scannable)

	if err != nil {
		return nil, err
	}

	err = residence.Associate(ctx, r, nil)

	if err != nil {
		return nil, err
	}

	return residence.toModel(nil), nil
}

func (r DBAisRepository) GetResidence(ctx context.Context, studentID int) (*models.Residence, error) {
	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", residenceFields, residenceTable), studentID)
	return makeResidenceModel(ctx, r, row)
}

func (r DBAisRepository) GetAllResidences(ctx context.Context) ([]*models.Residence, error) {
	errValue := []*models.Residence{}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s", residenceFields, residenceTable))

	if err != nil {
		return errValue, err
	}

	residences := []*models.Residence{}
	for rows.Next() {
		residence, err := makeResidenceModel(ctx, r, rows)

		if err != nil {
			return errValue, nil
		}

		residences = append(residences, residence)
	}

	return residences, nil
}
