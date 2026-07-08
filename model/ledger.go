package model

import "github.com/google/uuid"

type CreateResourceRequest struct {
	Resource []ItemRequest
}

type ItemRequest struct {
	Name     string `json:"name" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Unit     string `json:"unit" binding:"required"`
}

type CreateResourceResponse struct {
	LedgerID         uuid.UUID `json:"ledger_id"`
	DistributionCode string    `json:"distribution_code"`
	BlockNumber      string    `json:"block_number"`
}

type ConfirmResource struct {
	DistributionCode string `json:"distribution_code" binding:"required"`
}

type ConfirmResourceResponse struct {
	Resource []ItemRequest
}
