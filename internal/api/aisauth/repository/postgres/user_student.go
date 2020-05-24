package postgres

import (
	"context"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

// CREATE TABLE ais_user_student_binding(
//     id SERIAL,
//     user_id int NOT NULL references ais_user(id),
//     student_id int NOT NULL references Студент(id),
//     CONSTRAINT ais_student_binding_pk PRIMARY KEY (id)
// );

// type repoRole struct {
// 	ID  int
// 	Def string

// 	model *models.Role
// }

const studentTable = "Студент"
const userStudentBindingTable = "ais_user_student_binding"
const userStudentBindingFieldsNoID = "user_id,student_id"
const userStudentBindingStudentID = "student_id"
const userStudentBindingUserID = "user_id"

func (r AisAuthRepository) CreateUserStudentBinding(ctx context.Context, user *models.User, student *models.Student) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES ($1,$2)", userStudentBindingTable, userStudentBindingFieldsNoID),
		user.ID, student.ID,
	)

	return err
}

func (r AisAuthRepository) GetStudentUserID(ctx context.Context, student *models.Student) (int, error) {
	retVal := 0

	row, err := r.db.QueryContext(ctx, makeGetBindingQuery(userStudentBindingTable,
		userTable, "id", userStudentBindingUserID,
		studentTable, userStudentBindingStudentID,
	), student.ID)

	if err != nil {
		return retVal, err
	}

	if !row.Next() {
		return retVal, aisauth.ErrNoStudentUser
	}

	err = row.Scan(&retVal)

	if err != nil {
		return retVal, err
	}

	return retVal, err
}

func (r AisAuthRepository) GetUserStudentID(ctx context.Context, userID int) (*int, error) {
	var retVal *int

	row, err := r.db.QueryContext(ctx, makeGetBindingQuery(userStudentBindingTable,
		studentTable, "id", userStudentBindingStudentID,
		userTable, userStudentBindingUserID,
	), userID)

	if err != nil {
		return nil, err
	}

	if row.Next() {
		retVal = new(int)
		err = row.Scan(&retVal)
		if err != nil {
			return nil, err
		}
	}

	return retVal, err
}
