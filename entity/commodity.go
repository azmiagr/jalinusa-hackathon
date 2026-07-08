package entity

import "github.com/google/uuid"

type Commodity struct {
	CommodityID uuid.UUID `json:"commodity_id" gorm:"type:varchar(36);primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`

	LedgerItems []LedgerItem `json:"ledger_items" gorm:"foreignKey:CommodityID;references:CommodityID;constraint:onDelete:CASCADE"`
}
