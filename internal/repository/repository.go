package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository IUserRepository
	RoleRepository IRoleRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository: NewUserRepository(db),
		RoleRepository: NewRoleRepository(db),
	}
}
