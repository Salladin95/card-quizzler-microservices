package main

import (
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"log"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.AppCfg.RABBIT_URL)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	app := server.NewApp(cfg.AppCfg, rabbitConn)
	app.Start()
}
