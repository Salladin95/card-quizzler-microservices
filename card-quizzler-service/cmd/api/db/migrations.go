package migrations

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Folder{}, &models.Module{}, &models.Term{})
}
