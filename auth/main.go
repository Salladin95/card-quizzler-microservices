package main

import (
	"encoding/json"
	"fmt"
	"github.com/Salladin95/rmqtools"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type App struct {
	rabbit *amqp091.Connection
}

const (
	AmqpExchange = "broker"
	AmqpQueue    = "broker-queue"

	SignInKey = "auth.sign-in.command"
	SignUpKey = "auth.sign-up.command"
)

type SignInDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=6,required"`
}

type SighUpDto struct {
	Name     string `json:"name"  validate:"required,min=1"`
	Password string `json:"password"  validate:"required,min=6"`
	Email    string `json:"email"  validate:"required,email"`
	Birthday string `json:"birthday"  validate:"required,min=1"`
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://khalid:12345@localhost:5672/"
	}
	rabbitConn, err := rmqtools.ConnectToRabbit(rabbitURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app := App{rabbit: rabbitConn}

	consumer, err := rmqtools.NewConsumer(app.rabbit, AmqpExchange, AmqpQueue)

	if err != nil {
		log.Panic(err)
	}
	err = consumer.Listen([]string{SignInKey, SignUpKey}, handlePayload)
	if err != nil {
		log.Panic(err)
	}
}

func handlePayload(key string, payload []byte) {
	fmt.Print("START PROCESSING MESSAGE")
	switch key {
	case SignInKey:
		fmt.Printf("******* - %v\n\n", payload)
		var signInDto SignInDto
		err := json.Unmarshal(payload, &signInDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN IN *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %s\n", key)
		fmt.Printf("payload - %v\n\n", signInDto)
	case SignUpKey:
		var signUpDto SighUpDto
		err := json.Unmarshal(payload, &signUpDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN UP *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %v", key)
	default:
		log.Panic("handlePayload: unknown payload name")
	}
}
