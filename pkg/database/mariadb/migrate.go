package mariadb

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.Post{},
		&entity.AuditLog{},
		&entity.Device{},
		&entity.DeviceBinding{},
		&entity.LogisticLedger{},
		&entity.LedgerItem{},
		&entity.Distribution{},
	)

	if err != nil {
		return err
	}

	return nil
}
