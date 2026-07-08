package repository

import (
	"errors"

	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILogisticLedgerRepository interface {
	CreateLedger(tx *gorm.DB, ledger *entity.LogisticLedger) error
	GetLatestLedgerForUpdate(tx *gorm.DB, postID uuid.UUID) (*entity.LogisticLedger, error)
	GetLastLedger(tx *gorm.DB) (*entity.LogisticLedger, error)
	GetResourceListRequest(tx *gorm.DB) ([]*entity.LogisticLedger, error)
	GetLedger(tx *gorm.DB, param model.GetLedgerParam) (*entity.LogisticLedger, error)
}

type LogisticLedgerRepository struct {
	db *gorm.DB
}

func NewLogisticLedgerRepository(db *gorm.DB) ILogisticLedgerRepository {
	return &LogisticLedgerRepository{
		db: db,
	}
}

func (r *LogisticLedgerRepository) CreateLedger(tx *gorm.DB, ledger *entity.LogisticLedger) error {
	err := tx.Debug().Create(ledger).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *LogisticLedgerRepository) GetLatestLedgerForUpdate(tx *gorm.DB, postID uuid.UUID) (*entity.LogisticLedger, error) {
	var ledger entity.LogisticLedger

	err := tx.Debug().
		Where("post_id = ?", postID).
		Order("created_at DESC").
		First(&ledger).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ledger, nil
}

func (r *LogisticLedgerRepository) GetLastLedger(tx *gorm.DB) (*entity.LogisticLedger, error) {
	var ledger entity.LogisticLedger

	err := tx.
		Order("block_number DESC").
		First(&ledger).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &ledger, nil
}

func (r *LogisticLedgerRepository) GetResourceListRequest(tx *gorm.DB) ([]*entity.LogisticLedger, error) {
	var ledger []*entity.LogisticLedger
	err := tx.Debug().Find(&ledger).Error
	if err != nil {
		return nil, err
	}

	return ledger, nil
}

func (r *LogisticLedgerRepository) GetLedger(tx *gorm.DB, param model.GetLedgerParam) (*entity.LogisticLedger, error) {
	var ledger entity.LogisticLedger
	err := tx.Debug().Where(&param).Find(&ledger).Error
	if err != nil {
		return nil, err
	}

	return &ledger, nil
}
