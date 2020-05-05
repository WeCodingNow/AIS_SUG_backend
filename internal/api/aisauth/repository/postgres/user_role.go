package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

// CREATE TABLE ais_user_role(
//     id SERIAL,
//     def varchar(100) NOT NULL UNIQUE,
//     CONSTRAINT ais_user_role_pk PRIMARY KEY (id)
// );

// CREATE TABLE ais_role_binding(
//     id SERIAL,
//     ais_user_role_id int NOT NULL references ais_user_role(id) ON DELETE CASCADE,
//     ais_user_id int NOT NULL references ais_user(id) ON DELETE CASCADE,
//     CONSTRAINT ais_role_binding_pk PRIMARY KEY (id)
// );

// type repoRole struct {
// 	ID  int
// 	Def string

// 	model *models.Role
// }

// func (r repoRole) toModel() *models.Role {
// 	return &models.Role{
// 		ID:  r.ID,
// 		Def: r.Def,
// 	}
// }
const userTable = "ais_user"

const roleTable = "ais_user_role"
const userRoleBindingRoleID = "ais_user_role_id"
const userRoleBindingUserID = "ais_user_id"
const userRoleBindingTable = "ais_role_binding"
const userRoleBindingFieldsNoID = "ais_user_role_id,ais_user_id"

// const userRoleBindingFields = "id,ais_user_role_id,ais_user_id"

// const getBindingQuery = `
// SELECT %s.id FROM %s
// 	JOIN %s ON %s.%s = %s.id
// 	JOIN %s ON %s.%s = %s.id
// 	WHERE %s.id = $1
// `

// SELECT Студент.id FROM ais_user_student_binding
// 	  JOIN Студент ON ais_user_student_binding.student_id = Студент.id
// 	  JOIN ais_user ON ais_user_student_binding.user_id = ais_user.id

//rightTable - что получаем
func makeGetBindingQuery(bindingTable, leftTable, leftKey, rightTable, rightKey string) string {
	return fmt.Sprintf(
		`SELECT %s.id FROM %s
			JOIN %s ON %s.%s = %s.id
			JOIN %s ON %s.%s = %s.id
			WHERE %s.id = $1`,
		leftTable, bindingTable,
		leftTable, bindingTable, leftKey, leftTable,
		rightTable, bindingTable, rightKey, rightTable,
		rightTable,
	)
}

func (r AisAuthRepository) CreateUserRoleBinding(ctx context.Context, user *models.User, role *models.Role) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES ($1,$2)", userRoleBindingTable, userRoleBindingFieldsNoID),
		role.ID, user.ID,
	)

	return err
}

func (r AisAuthRepository) GetUserRoleID(ctx context.Context, user *models.User) (int, error) {
	retVal := 0

	query := makeGetBindingQuery(userRoleBindingTable,
		roleTable, userRoleBindingRoleID,
		userTable, userRoleBindingUserID,
	)

	log.Print(query)

	row := r.db.QueryRowContext(ctx, query, user.ID)

	err := row.Scan(&retVal)
	return retVal, err
}
