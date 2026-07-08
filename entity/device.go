package entity

import "github.com/google/uuid"

type Device struct {
	DeviceID   uuid.UUID `json:"device_id" gorm:"type:varchar(36);primaryKey"`
	DeviceName string    `json:"device_name" gorm:"type:varchar(255)"`

	DeviceBindings []DeviceBinding `json:"device_bindings" gorm:"foreignKey:DeviceID;references:DeviceID;constraint:onDelete:CASCADE"`
}
