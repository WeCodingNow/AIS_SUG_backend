package postgres

// import (
// 	"context"
// 	"log"

// 	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
// )

// type Role struct {
// 	ID  int    `json:"id"`
// 	Def string `json:"def"`
// }

// func toPostgresRole(r *models.Role) *Role {
// 	return &Role{
// 		ID: r.ID,
// 	}
// }

// func toModelRole(r *Role) *models.Role {
// 	return &models.Role{
// 		ID: r.ID,
// 	}
// }

// const getUserRoleQuery = `SELECT ucl.id, ucl.def as user_class_id
// 	FROM ais_role_binding as rb
// 	JOIN ais_user as u ON rb.ais_user_id = u.id
// 	JOIN ais_user_class as ucl ON rb.ais_user_class_id = ucl.id
// 	WHERE ais_user_id = $1;`

// func (r UserRepository) GetUserRole(ctx context.Context, user *models.User) (*models.Role, error) {
// 	user, err := r.GetUserByName(ctx, user.Username)

// 	if err != nil {
// 		return nil, err
// 	}

// 	row := r.db.QueryRowContext(
// 		ctx, getUserRoleQuery,
// 		user.ID,
// 	)

// 	role := new(Role)
// 	if err := row.Scan(&role.ID, &role.Def); err != nil {
// 		log.Printf("PG err: %+v", err)
// 	}

// 	return toModelRole(role), nil
// }
