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

type ResourceRequestList struct {
	Resources []ResourceResponse
}

type ResourceResponse struct {
	LedgerID           uuid.UUID `json:"ledger_id"`
	PostName           string    `json:"post_name"`
	DistributionCode   string    `json:"distribution_code"`
	DistributionStatus string    `json:"distribution_status"`
	BlockNumber        string    `json:"block_number"`
	Items              []ItemRequest
}

type GetResourceDetail struct {
	Items      []ItemRequest
	Status     string `json:"status"`
	HashLedger string `json:"hash_ledger"`
}

type GetLedgerParam struct {
	LedgerID uuid.UUID `json:"ledger_id"`
}
