package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLogResponse struct {
	AuditID   uuid.UUID  `json:"audit_id"`
	UserID    *uuid.UUID `json:"user_id"`
	Action    string     `json:"action"`
	CreatedAt time.Time  `json:"created_at"`
}
