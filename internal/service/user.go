package service

import (
	"errors"

	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/azmiagr/jalinusa-hackathon/pkg/bcrypt"
	"github.com/azmiagr/jalinusa-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/jalinusa-hackathon/pkg/errors"
	"github.com/azmiagr/jalinusa-hackathon/pkg/jwt"

	"github.com/azmiagr/jalinusa-hackathon/internal/repository"
	"github.com/azmiagr/jalinusa-hackathon/model"

	"gorm.io/gorm"
)

type IUserService interface {
	Login(param model.UserLoginRequest) (*model.UserLoginResponse, error)
	GetUser(param model.GetUserParam) (*entity.User, error)
}

type UserService struct {
	db       *gorm.DB
	userRepo repository.IUserRepository
	roleRepo repository.IRoleRepository
	bcrypt   bcrypt.Interface
	jwtAuth  jwt.Interface
}

func NewUserService(userRepo repository.IUserRepository, roleRepo repository.IRoleRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		db:       mariadb.Connection,
		userRepo: userRepo,
		roleRepo: roleRepo,
		bcrypt:   bcrypt,
		jwtAuth:  jwtAuth,
	}
}

func (s *UserService) Login(param model.UserLoginRequest) (*model.UserLoginResponse, error) {
	user, err := s.userRepo.GetUser(s.db, model.GetUserParam{
		Username: param.Username,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Unauthorized("invalid email or password")
		}
		return nil, apperrors.InternalServer("failed to get user")
	}

	err = s.bcrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid email or password")
	}

	if user.Status != "active" {
		return nil, apperrors.Unauthorized("account is not active")
	}

	role, err := s.roleRepo.GetRoleByID(s.db, user.RoleID)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get role")
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID, role.RoleName)
	if err != nil {
		return nil, apperrors.InternalServer("failed to generate token")
	}

	return &model.UserLoginResponse{
		Token: token,
	}, nil

}

func (s *UserService) GetUser(param model.GetUserParam) (*entity.User, error) {
	return s.userRepo.GetUser(s.db, param)
}
