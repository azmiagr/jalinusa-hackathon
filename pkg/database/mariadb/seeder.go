package mariadb

import (
	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) error {
	roles := []entity.Role{
		{
			RoleID:   uuid.MustParse("10d0dea6-7a7c-11f1-a822-66de32dbe509"),
			RoleName: "admin",
		},
	}

	for _, role := range roles {
		err := db.
			Where("role_name = ?", role.RoleName).
			FirstOrCreate(&role).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func SeedUser(db *gorm.DB) error {
	password, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	var role entity.Role

	err = db.
		Where("role_name = ?", "admin").
		First(&role).Error
	if err != nil {
		return err
	}

	user := entity.User{
		UserID:   uuid.New(),
		RoleID:   role.RoleID,
		Username: "admin",
		Password: string(password),
	}

	err = db.
		Where("username = ?", user.Username).
		FirstOrCreate(&user).Error

	return err
}
