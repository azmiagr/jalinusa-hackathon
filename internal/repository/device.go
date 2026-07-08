package repository

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"
	"gorm.io/gorm"
)

type IDeviceRepository interface {
	FindOrCreateDevice(tx *gorm.DB, device *entity.Device) (*entity.Device, error)
	CreateDeviceBinding(tx *gorm.DB, binding *entity.DeviceBinding) error
}

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) IDeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) FindOrCreateDevice(tx *gorm.DB, device *entity.Device) (*entity.Device, error) {
	var result entity.Device

	err := tx.
		Where("device_name = ?", device.DeviceName).
		FirstOrCreate(&result, device).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *DeviceRepository) CreateDeviceBinding(tx *gorm.DB, binding *entity.DeviceBinding) error {
	err := tx.Debug().Create(binding).Error
	if err != nil {
		return err
	}

	return nil
}
