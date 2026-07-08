package entity

import "github.com/google/uuid"

type Distribution struct {
	DistributionID uuid.UUID `json:"distribution_id" gorm:"type:varchar(36);primaryKey"`
	LedgerID       uuid.UUID `json:"ledger_id" gorm:"type:varchar(36)"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	Status         string    `json:"status" gorm:"type:enum('diajukan', 'diproses', 'pengiriman', 'terdistribusi', 'disetujui')"`
}
