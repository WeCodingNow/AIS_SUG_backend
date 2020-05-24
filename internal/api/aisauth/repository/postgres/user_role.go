package postgres

import (
	"context"
	"fmt"
	"strings"

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

type repoRole struct {
	ID  int
	Def string

	model *models.Role
}

func (r repoRole) toModel() *models.Role {
	return &models.Role{
		ID:  r.ID,
		Def: r.Def,
	}
}

func NewRepoRole() *repoRole {
	return &repoRole{}
}

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
func makeGetBindingQueryID(bindingTable, leftTable, leftKey, rightTable, rightKey string) string {
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

func prefixFields(table, fieldsString string) string {
	fields := strings.Split(fieldsString, ",")

	for i := range fields {
		fields[i] = fmt.Sprintf("%s.%s", table, fields[i])
	}

	return strings.Join(fields, ",")
}

//rightTable - что получаем
func makeGetBindingQuery(bindingTable, leftTable, leftTableFields, leftKey, rightTable, rightKey string) string {
	return fmt.Sprintf(
		`SELECT %s FROM %s
			JOIN %s ON %s.%s = %s.id
			JOIN %s ON %s.%s = %s.id
			WHERE %s.id = $1`,
		prefixFields(leftTable, leftTableFields), bindingTable,
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
	userRoleID := 0

	query := makeGetBindingQueryID(userRoleBindingTable,
		roleTable, userRoleBindingRoleID,
		userTable, userRoleBindingUserID,
	)

	row := r.db.QueryRowContext(ctx, query, user.ID)

	err := row.Scan(&userRoleID)
	return userRoleID, err
}

func (r AisAuthRepository) GetUserRole(ctx context.Context, userID int) (*models.Role, error) {
	repoRole := NewRepoRole()

	query := makeGetBindingQuery(userRoleBindingTable,
		roleTable, "id,def", userRoleBindingRoleID,
		userTable, userRoleBindingUserID,
	)

	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&repoRole.ID, &repoRole.Def)

	return repoRole.toModel(), err
}

func (r AisAuthRepository) GetRoles(ctx context.Context) ([]*models.Role, error) {
	roles := make([]*models.Role, 0)

	query := "select id, def from ais_user_role"

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		role := new(models.Role)

		err = rows.Scan(&role.ID, &role.Def)

		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r AisAuthRepository) PromoteUser(ctx context.Context, userID int, roleID int) error {
	query := `UPDATE ais_role_binding
	SET ais_user_role_id = $2
	WHERE ais_user_id = $1;
	`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)

	return err
}
