package migrations

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}, &models.Folder{}, &models.Module{}, &models.Term{}); err != nil {
		log.Fatal("failed to apply migrations")
	}
}
