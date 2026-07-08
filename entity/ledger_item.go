package entity

import "github.com/google/uuid"

type LedgerItem struct {
	LedgerItemID uuid.UUID `json:"ledger_item_id" gorm:"type:varchar(36);primaryKey"`
	LedgerID     uuid.UUID `json:"ledger_id" gorm:"type:varchar(36)"`
	CommodityID  uuid.UUID `json:"commodity_id" gorm:"type:varchar(36)"`
	Quantity     int       `json:"quantity" gorm:"type:int"`
	Unit         string    `json:"unit" gorm:"type:varchar(50)"`
}
