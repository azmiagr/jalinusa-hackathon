package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AuditID   uuid.UUID `json:"audit_id" gorm:"type:varchar(36);primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	Action    string    `json:"action" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
