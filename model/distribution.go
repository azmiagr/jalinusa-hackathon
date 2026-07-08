package model

import "github.com/google/uuid"

type DistributionParam struct {
	Code     string    `json:"-"`
	LedgerID uuid.UUID `json:"-"`
}
