package main

import (
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/server"
	"github.com/Salladin95/rmqtools"
	"log"
	"os"
)

// main is the entry point of the application.
func main() {
	// Load application configuration.
	cfg, err := config.NewConfig()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.AppCfg.RABBIT_URL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Ensure the RabbitMQ connection is closed when the main function exits.
	defer rabbitConn.Close()

	// Create a new instance of the application using the loaded configuration and RabbitMQ connection & start it
	server.NewApp(cfg.AppCfg, rabbitConn).Start()
}

//consumer, err := rmqtools.NewConsumer(app.rabbit, AmqpExchange, AmqpQueue)
//if err != nil {
//	log.Panic(err)
//}
//err = consumer.Listen([]string{SignInKey, SignUpKey}, handlePayload)
//if err != nil {
//	log.Panic(err)
//}
