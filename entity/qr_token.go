package entity

import "github.com/google/uuid"

type QRToken struct {
	QRTokenID uuid.UUID `json:"qr_token_id" gorm:"type:varchar(36);primaryKey"`
	LedgerID  uuid.UUID `json:"ledger_id" gorm:"type:varchar(36)"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
}
