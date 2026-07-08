package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeviceBinding struct {
	DeviceBindingID uuid.UUID `json:"device_binding_id" gorm:"type:varchar(36);primaryKey"`
	DeviceID        uuid.UUID `json:"device_id" gorm:"type:varchar(36)"`
	PostID          uuid.UUID `json:"post_id" gorm:"type:varchar(36)"`
	Status          string    `json:"status" gorm:"type:enum('success','failed')"`
	BoundBy         uuid.UUID `json:"bound_by" gorm:"type:varchar(36)"`
	BoundAt         time.Time `json:"bound_at" gorm:"type:datetime"`
}
