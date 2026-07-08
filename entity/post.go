package entity

import "github.com/google/uuid"

type Post struct {
	PostID     uuid.UUID `json:"post_id" gorm:"type:varchar(36);primaryKey"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:varchar(36);"`
	PostCode   string    `json:"post_code" gorm:"type:varchar(50);not null;uniqueIndex"`
	PostName   string    `json:"post_name" gorm:"type:varchar(255)"`
	Coordinate string    `json:"coordinate" gorm:"type:varchar(255)"`
	Capacity   int       `json:"capacity" gorm:"type:int"`
	QRCode     string    `json:"qr_code" gorm:"type:varchar(255)"`
	Status     string    `json:"status" gorm:"type:enum('active', 'inactive');default:'active'"`

	LogisticLedgers []LogisticLedger `json:"logistic_ledgers" gorm:"foreignKey:PostID;references:PostID;constraint:onDelete:CASCADE"`
	DeviceBindings  []DeviceBinding  `json:"device_bindings" gorm:"foreignKey:PostID;references:PostID;constraint:onDelete:CASCADE"`
}
