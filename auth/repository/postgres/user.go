package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/WeCodingNow/AIS_SUG_backend/models"
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

func toModel(u *User) *models.User {
	return &models.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	model := toPostgresUser(user)

	row := r.db.QueryRowContext(
		ctx, `INSERT INTO ais_user(username,password) VALUES ( $1, $2 ) RETURNING id;`,
		model.Username, model.Password,
	)

	var newId int
	if err := row.Scan(&newId); err != nil {
		return err
	} else {
		user.ID = newId
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
		return nil, err
	}

	return toModel(user), nil
}

func (r UserRepository) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	row := r.db.QueryRowContext(
		ctx, `SELECT * FROM ais_user WHERE username = $1`,
		username,
	)

	user := new(User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		log.Printf("PG err: %+v", err)
		return nil, err
	}

	return toModel(user), nil
}
