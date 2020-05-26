package usecase

import (
	"context"
	"log"

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

func (aue AisAuthUseCase) GetUserStudentID(ctx context.Context, userID int) (*int, error) {
	return aue.repo.GetUserStudentID(ctx, userID)
}

func (aue AisAuthUseCase) GetStudentUserID(ctx context.Context, student *models.Student) (int, error) {
	return aue.repo.GetStudentUserID(ctx, student)
}

func isSameStudent(lhs, rhs *models.Student) bool {
	retBool := false

	if lhs.Name == rhs.Name &&
		lhs.SecondName == rhs.SecondName &&
		lhs.Group.ID == rhs.Group.ID {
		log.Printf("now sravnivaem otchestva")
		if (lhs.ThirdName != nil && rhs.ThirdName != nil && *lhs.ThirdName == *rhs.ThirdName) ||
			(lhs.ThirdName == nil && rhs.ThirdName == nil) {
			retBool = true
		}
	}

	return retBool
}

func (aue AisAuthUseCase) AssignOrCreateStudentWithCreds(
	ctx context.Context,
	user *models.User,
	targetStudent *models.Student,
) error {
	students, err := aue.aisUC.GetAllStudents(ctx)

	if err != nil {
		return err
	}

	for _, student := range students {
		if isSameStudent(targetStudent, student) {
			foundUser, err := aue.GetStudentUserID(ctx, student)
			if err != nil {
				if err != aisauth.ErrNoStudentUser {
					return err
				}
			}

			if foundUser != 0 {
				return aisauth.ErrDuplicateUser
			}

			targetStudent.ID = student.ID
		}
	}

	return aue.CreateStudentWithCreds(ctx, user, &models.Role{ID: 3}, targetStudent)
}

func (aue AisAuthUseCase) CreateStudentWithCreds(
	ctx context.Context,
	user *models.User, role *models.Role, student *models.Student,
) error {
	var err error
	if user.ID == 0 {
		newUser, err := aue.authUC.CreateUser(ctx, user.Username, user.Password)

		if err != nil {
			return err
		}

		user = newUser
		err = aue.repo.CreateUserRoleBinding(ctx, user, role)

		if err != nil {
			return err
		}
	}

	if student != nil {
		_, err = aue.aisUC.GetStudent(ctx, student.ID)
		if err != nil {
			if err == ais.ErrStudentNotFound {
				student, err = aue.aisUC.CreateStudent(ctx, student.Name, student.SecondName, student.ThirdName, student.Group.ID, student.Residence.ID)

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

func (aue AisAuthUseCase) GetStudentsWithUsers(ctx context.Context) ([]aisauth.StudentWithUserAndRole, error) {
	students, err := aue.aisUC.GetAllStudents(ctx)

	if err != nil {
		return nil, err
	}

	retBindings := make([]aisauth.StudentWithUserAndRole, len(students))

	for i, st := range students {
		retBindings[i].StudentID = st.ID
		userID, err := aue.repo.GetStudentUserID(ctx, st)

		if err != nil {
			if err == aisauth.ErrNoStudentUser {
				continue
			}
			return nil, err
		}

		if userID != 0 {
			retBindings[i].UserID = new(int)
			*retBindings[i].UserID = userID

			roleID, err := aue.repo.GetUserRoleID(ctx, &models.User{ID: userID})

			if err != nil {
				return nil, err
			}

			retBindings[i].RoleID = new(int)
			*retBindings[i].RoleID = roleID
		}
	}

	return retBindings, nil
}

func (aue AisAuthUseCase) GetRoles(ctx context.Context) ([]*models.Role, error) {
	return aue.repo.GetRoles(ctx)
}

func (aue AisAuthUseCase) PromoteUser(ctx context.Context, userID int, roleID int) error {
	return aue.repo.PromoteUser(ctx, userID, roleID)
}

func (aue AisAuthUseCase) GetInfo(ctx context.Context, userID int) (*aisauth.StudentInfo, error) {
	info := &aisauth.StudentInfo{UserID: userID}
	var err error

	info.StudentID, err = aue.GetUserStudentID(ctx, userID)

	if err != nil {
		return nil, err
	}

	if info.StudentID != nil {
		student, err := aue.aisUC.GetStudent(ctx, *info.StudentID)

		if err != nil {
			return nil, err
		}

		info.GroupID = new(int)
		*info.GroupID = student.Group.ID
	}

	log.Print(*info.GroupID)
	return info, nil
}
