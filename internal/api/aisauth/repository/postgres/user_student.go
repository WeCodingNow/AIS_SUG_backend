package postgres

import (
	"context"
	"fmt"

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

func (r AisAuthRepository) GetUserStudentID(ctx context.Context, user *models.User) (int, error) {
	retVal := 0

	row := r.db.QueryRowContext(ctx, makeGetBindingQuery(userStudentBindingTable,
		studentTable, userStudentBindingStudentID,
		userTable, userStudentBindingUserID,
	))

	err := row.Scan(&retVal)
	return retVal, err
}
