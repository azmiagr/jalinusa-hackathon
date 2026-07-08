package service

import (
	"github.com/azmiagr/jalinusa-hackathon/internal/repository"
	"github.com/azmiagr/jalinusa-hackathon/pkg/bcrypt"
	"github.com/azmiagr/jalinusa-hackathon/pkg/jwt"
)

type Service struct {
	UserService   IUserService
	PostService   IPostService
	LedgerService ILedgerService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
	return &Service{
		UserService:   NewUserService(repository.UserRepository, repository.RoleRepository, bcrypt, jwtAuth),
		PostService:   NewPostService(repository.PostRepository, repository.UserRepository, repository.DeviceRepository),
		LedgerService: NewLedgerService(repository.LogisticLedgerRepository, repository.DistributionRepository, repository.LedgerItemRepository, repository.UserRepository, repository.PostRepository, repository.AuditLogRepository),
	}
}
