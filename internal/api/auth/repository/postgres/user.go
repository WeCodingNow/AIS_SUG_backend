package postgres

import (
	"context"
	"log"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func toPostgresUser(u *models.User) *User {
	return &User{
		Username: u.Username,
		Password: u.Password,
	}
}

func toModelUser(u *User) *models.User {
	return &models.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	userModel := toPostgresUser(user)
	// roleModel := toPostgresRole(role)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx,
		`INSERT INTO ais_user(username,password) VALUES ( $1, $2 ) RETURNING id;`,
		userModel.Username, userModel.Password,
	)

	var newID int
	if err := row.Scan(&newID); err != nil {
		tx.Rollback()
		return err
	} else {
		user.ID = newID
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	// log.Print(roleModel, userModel)

	// _, err = tx.ExecContext(ctx,
	// 	`INSERT INTO ais_role_binding(ais_user_class_id, ais_user_id, user_class_confirmed) VALUES ( $1, $2, true)`,
	// 	roleModel.ID, newID,
	// )

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	row := r.db.QueryRowContext(
		ctx, `SELECT * FROM ais_user WHERE username = $1 AND password = $2`,
		username, password,
	)

	user := new(User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		log.Printf("PG err: %+v", err)
		return nil, auth.ErrUserNotFound
	}

	return toModelUser(user), nil
}

func (r UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT * FROM ais_user WHERE id = $1`, id)

	user := new(User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		log.Printf("PG err: %+v", err)
		return nil, err
	}

	return toModelUser(user), nil
}
