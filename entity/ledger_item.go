package entity

import "github.com/google/uuid"

type LedgerItem struct {
	LedgerItemID uuid.UUID `json:"ledger_item_id" gorm:"type:varchar(36);primaryKey"`
	LedgerID     uuid.UUID `json:"ledger_id" gorm:"type:varchar(36)"`
	Name         string    `json:"name" gorm:"type:varchar(100)"`
	Quantity     int       `json:"quantity" gorm:"type:int"`
	Unit         string    `json:"unit" gorm:"type:varchar(50)"`
}
