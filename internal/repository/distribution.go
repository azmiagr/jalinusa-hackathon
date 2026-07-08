package repository

import (
	"errors"

	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDistributionRepository interface {
	CreateDistribution(tx *gorm.DB, distribution *entity.Distribution) error
	GetLastDistribution(tx *gorm.DB) (*entity.Distribution, error)
	GetDistribution(tx *gorm.DB, param model.DistributionParam) (*entity.Distribution, error)
	UpdateDistribution(tx *gorm.DB, distribution *entity.Distribution) error
	GetDistributionsByLedgerID(tx *gorm.DB, ledgerID uuid.UUID) ([]*entity.Distribution, error)
	GetSubmittedDistribution(tx *gorm.DB) (int64, error)
	GetDeliveredDistribution(tx *gorm.DB) (int64, error)
	GetAcceptedDistribution(tx *gorm.DB) (int64, error)
	GetUnfinishedDistribution(tx *gorm.DB) (int64, error)
	GetAllDistributionCount(tx *gorm.DB) (int64, error)
}

type DistributionRepository struct {
	db *gorm.DB
}

func NewDistributionRepository(db *gorm.DB) IDistributionRepository {
	return &DistributionRepository{
		db: db,
	}
}

func (r *DistributionRepository) CreateDistribution(tx *gorm.DB, distribution *entity.Distribution) error {
	err := tx.Debug().Create(distribution).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DistributionRepository) GetLastDistribution(tx *gorm.DB) (*entity.Distribution, error) {
	var distribution entity.Distribution

	err := tx.
		Order("code DESC").
		First(&distribution).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &distribution, nil
}

func (r *DistributionRepository) GetDistribution(tx *gorm.DB, param model.DistributionParam) (*entity.Distribution, error) {
	var distribution entity.Distribution

	err := tx.Debug().Where(&param).First(&distribution).Error
	if err != nil {
		return nil, err
	}

	return &distribution, nil

}

func (r *DistributionRepository) UpdateDistribution(tx *gorm.DB, distribution *entity.Distribution) error {
	err := tx.Debug().Save(distribution).Error
	if err != nil {
		return err
	}

	return err
}

func (r *DistributionRepository) GetDistributionsByLedgerID(tx *gorm.DB, ledgerID uuid.UUID) ([]*entity.Distribution, error) {
	var distribution []*entity.Distribution
	err := tx.Debug().Where("ledger_id = ?", ledgerID).Find(&distribution).Error
	if err != nil {
		return nil, err
	}

	return distribution, nil
}

func (r *DistributionRepository) GetSubmittedDistribution(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.Distribution{}).Where("status = 'diajukan' AND created_at <= INTERVAL 1 HOUR + NOW()").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DistributionRepository) GetDeliveredDistribution(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.Distribution{}).Where("status = 'pengiriman'").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DistributionRepository) GetAcceptedDistribution(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.Distribution{}).Where("status = 'diterima'").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DistributionRepository) GetUnfinishedDistribution(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.Distribution{}).Where("status != 'diterima'").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *DistributionRepository) GetAllDistributionCount(tx *gorm.DB) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.Distribution{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
