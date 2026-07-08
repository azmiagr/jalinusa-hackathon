package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository           IUserRepository
	RoleRepository           IRoleRepository
	PostRepository           IPostRepository
	DeviceRepository         IDeviceRepository
	LogisticLedgerRepository ILogisticLedgerRepository
	DistributionRepository   IDistributionRepository
	LedgerItemRepository     ILedgerItemRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:           NewUserRepository(db),
		RoleRepository:           NewRoleRepository(db),
		PostRepository:           NewPostRepository(db),
		DeviceRepository:         NewDeviceRepository(db),
		LogisticLedgerRepository: NewLogisticLedgerRepository(db),
		DistributionRepository:   NewDistributionRepository(db),
		LedgerItemRepository:     NewLedgerItemRepository(db),
	}
}
