package main

const (
	AmqpExchange = "broker"
	AmqpQueue    = "broker-queue"

	SignInKey        = "auth.sign-in.command"
	SignUpKey        = "auth.sign-up.command"
	defaultRabbitURL = "amqp://khalid:12345@localhost:5672" //local
	//defaultRabbitURL = "amqp://khalid:12345@rabbitmq:5672" // from container
	gRPCPort = "8090"
)
