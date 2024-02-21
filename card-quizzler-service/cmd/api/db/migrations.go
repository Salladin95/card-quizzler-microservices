package migrations

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	//db.AutoMigrate(&models.Folder{}, &models.Module{}, &models.Term{}, &models.User{})
}
