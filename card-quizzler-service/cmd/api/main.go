package main

import (
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/db/migrations"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/server"
	lib "github.com/Salladin95/card-quizzler-microservices/shared"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	s := echo.New()
	cfg, err := config.NewConfig()

	if err != nil {
		s.Logger.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	migrations.Migrate(db)

	services := lib.InitializeServices(cfg.ServicesCfg)
	// Close connections when main function exits
	defer services.Rabbit.Close() // Close RabbitMQ connection
	defer services.Redis.Close()  // Close Redis connection

	server.NewApp(cfg, services.Rabbit, db, services.Redis).Start()
}
