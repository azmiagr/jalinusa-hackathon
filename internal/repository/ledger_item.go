package repository

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILedgerItemRepository interface {
	CreateLedgerItem(tx *gorm.DB, item *entity.LedgerItem) error
	GetLedgerItemByLedgerID(tx *gorm.DB, ledgerID uuid.UUID) ([]*entity.LedgerItem, error)
}

type LedgerItemRepository struct {
	db *gorm.DB
}

func NewLedgerItemRepository(db *gorm.DB) ILedgerItemRepository {
	return &LedgerItemRepository{
		db: db,
	}
}

func (r *LedgerItemRepository) CreateLedgerItem(tx *gorm.DB, item *entity.LedgerItem) error {
	err := tx.Debug().Create(item).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *LedgerItemRepository) GetLedgerItemByLedgerID(tx *gorm.DB, ledgerID uuid.UUID) ([]*entity.LedgerItem, error) {
	var items []*entity.LedgerItem
	err := tx.Debug().Where("ledger_id = ?", ledgerID).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
