package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository   IUserRepository
	RoleRepository   IRoleRepository
	PostRepository   IPostRepository
	DeviceRepository IDeviceRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:   NewUserRepository(db),
		RoleRepository:   NewRoleRepository(db),
		PostRepository:   NewPostRepository(db),
		DeviceRepository: NewDeviceRepository(db),
	}
}
