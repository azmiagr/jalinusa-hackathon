package repository

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetRoleByID(tx *gorm.DB, roleID uuid.UUID) (*entity.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) GetRoleByID(tx *gorm.DB, roleID uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := tx.Debug().Where("role_id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}
