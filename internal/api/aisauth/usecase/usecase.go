package usecase

import (
	"context"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/aisauth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/auth"
	"github.com/WeCodingNow/AIS_SUG_backend/internal/api/models"
)

type AisAuthUseCase struct {
	repo   aisauth.AisAuthRepository
	aisUC  ais.UseCase
	authUC auth.UseCase
}

func NewAisAuthUseCase(
	aisauthRepo aisauth.AisAuthRepository,
	aisUC ais.UseCase,
	authUC auth.UseCase,
) *AisAuthUseCase {
	return &AisAuthUseCase{
		aisauthRepo,
		aisUC,
		authUC,
	}
}

func (aue AisAuthUseCase) GetUserRoleID(ctx context.Context, user *models.User) (int, error) {
	return aue.repo.GetUserRoleID(ctx, user)
}

func (aue AisAuthUseCase) GetUserRole(ctx context.Context, userID int) (*models.Role, error) {
	return aue.repo.GetUserRole(ctx, userID)
}

func (aue AisAuthUseCase) GetUserStudentID(ctx context.Context, user *models.User) (int, error) {
	return aue.repo.GetUserStudentID(ctx, user)
}

func (aue AisAuthUseCase) CreateStudentWithCreds(
	ctx context.Context,
	user *models.User, role *models.Role, student *models.Student,
) error {
	user, err := aue.authUC.CreateUser(ctx, user.Username, user.Password)

	if err != nil {
		return err
	}

	err = aue.repo.CreateUserRoleBinding(ctx, user, role)

	if err != nil {
		return err
	}

	if student != nil {
		student, err = aue.aisUC.GetStudent(ctx, student.ID)
		if err != nil {
			if err == ais.ErrStudentNotFound {
				student, err = aue.aisUC.CreateStudent(ctx, student.Name, student.SecondName, student.ThirdName, student.Group.ID)

				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if err != nil {
			return err
		}
		err = aue.repo.CreateUserStudentBinding(ctx, user, student)

		if err != nil {
			return err
		}
	}

	return nil
}
