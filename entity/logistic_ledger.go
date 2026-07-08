package entity

import (
	"time"

	"github.com/google/uuid"
)

type LogisticLedger struct {
	LedgerID    uuid.UUID `json:"ledger_id" gorm:"type:varchar(36);primaryKey"`
	PostID      uuid.UUID `json:"post_id" gorm:"type:varchar(36)"`
	PrevHash    string    `json:"prev_hash" gorm:"type:varchar(255)"`
	CurrentHash string    `json:"current_hash" gorm:"type:varchar(255)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	BlockNumber string    `json:"block_number" gorm:"varchar(10)"`

	LedgerItems   []LedgerItem   `json:"ledger_items" gorm:"foreignKey:LedgerID;references:LedgerID;constraint:onDelete:CASCADE"`
	Distributions []Distribution `json:"distributions" gorm:"foreignKey:LedgerID;references:LedgerID;constraint:onDelete:CASCADE"`
}
