package repository

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"
	"gorm.io/gorm"
)

type IAuditLogRepository interface {
	CreateAuditLog(tx *gorm.DB, audit *entity.AuditLog) error
	GetAllAuditLog(tx *gorm.DB) ([]*entity.AuditLog, error)
}

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) IAuditLogRepository {
	return &AuditLogRepository{
		db: db,
	}
}

func (r *AuditLogRepository) CreateAuditLog(tx *gorm.DB, audit *entity.AuditLog) error {
	err := tx.Debug().Create(audit).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AuditLogRepository) GetAllAuditLog(tx *gorm.DB) ([]*entity.AuditLog, error) {
	var audit []*entity.AuditLog
	err := tx.Debug().Find(&audit).Error
	if err != nil {
		return nil, err
	}

	return audit, nil
}
