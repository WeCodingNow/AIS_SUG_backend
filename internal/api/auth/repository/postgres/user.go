package postgres

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
	"github.com/WeCodingNow/AIS_SUG_backend/pkg/pgorm"
)

type repoUser struct {
	ID       int
	Username string
	Password string
}

func NewRepoUser() *repoUser {
	return &repoUser{}
}

func (u *repoUser) Fill(scannable pgorm.Scannable) {
	scannable.Scan(&u.ID, &u.Username, &u.Password)
}

func (u *repoUser) GetID() int {
	return u.ID
}

const userTable = "ais_user"
const userFields = "id,username,password"

func (r repoUser) GetDescription() pgorm.ModelDescription {
	return pgorm.ModelDescription{
		Table:  userTable,
		Fields: userFields,
	}
}

func (u *repoUser) toModel() *models.User {
	return &models.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

func (u *repoUser) AcceptDep(interface{}) error {
	return nil
}

func (r UserRepository) CreateUser(ctx context.Context, username, password string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO ais_user(username,password) VALUES ( $1, $2 ) RETURNING id`,
		username, password,
	)

	var newID int
	err := row.Scan(&newID)
	if err != nil {
		return nil, err
	}

	user := &repoUser{
		ID:       newID,
		Username: username,
		Password: password,
	}

	return user.toModel(), nil
}

func (r UserRepository) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	user := NewRepoUser()
	filler, err := pgorm.MakeFillerWithFilter(
		ctx, r.db, userFields, userTable,
		pgorm.FilteredFields{"username", "password"},
		username, password,
	)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, auth.ErrUserNotFound
	}

	err = filler.Fill(user)

	return user.toModel(), nil
}

func (r UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := NewRepoUser()
	filler, err := pgorm.MakeFiller(ctx, r.db, userFields, userTable, &id)

	if err != nil {
		return nil, err
	}

	if !filler.Next() {
		return nil, auth.ErrUserNotFound
	}

	err = filler.Fill(user)

	return user.toModel(), nil
}
