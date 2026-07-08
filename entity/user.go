package entity

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	RoleID   uuid.UUID `json:"role_id" gorm:"type:varchar(36)"`
	Username string    `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Password string    `json:"password" gorm:"type:varchar(255);not null"`
	Status   string    `json:"status" gorm:"type:enum('active','inactive');default:'inactive'"`

	Posts          []Post          `json:"posts" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	AuditLogs      []AuditLog      `json:"audit_logs" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	DeviceBindings []DeviceBinding `json:"device_bindings" gorm:"foreignKey:BoundBy;references:UserID;constraint:onDelete:CASCADE"`
	Distributions  []Distribution  `json:"distributions" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
}
