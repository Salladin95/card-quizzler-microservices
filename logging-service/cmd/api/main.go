package main

import (
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/logging-service/cmd/api/entities"
	"github.com/Salladin95/rmqtools"
	"log"
	"os"
)

func main() {
	// Load application configuration.
	cfg, err := config.NewConfig()

	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	consumer, err := rmqtools.NewConsumer(
		rabbitConn,
		constants.AmqpExchange,
		constants.AmqpQueue,
	)
	if err != nil {
		// TODO: REPLACE
		fmt.Println(err)
	}
	err = consumer.Listen(
		[]string{
			constants.LogCommand,
		},
		func(routingKey string, payload []byte) {
			var data entities.LogMessage
			err := json.Unmarshal(payload, &data)
			if err != nil {
				fmt.Println("failed to unmarshall data")
				return
			}
			fmt.Printf("[logging service][routingkey - %s] message - %v\n", routingKey, data.Message)
			// TODO: SAVE LOG IN A SOURCE
		},
	)
	if err != nil {
		// TODO: REPLACE
		fmt.Println(err)
	}
}
