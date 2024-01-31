package main

import (
	"context"
	"fmt"
	auth "github.com/Salladin95/card-quizzler-microservices/auth-service/proto"
	"github.com/Salladin95/rmqtools"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type App struct {
	rabbit *amqp091.Connection
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = defaultRabbitURL
	}

	rabbitConn, err := rmqtools.ConnectToRabbit(rabbitURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := App{rabbit: rabbitConn}
	go app.gRPCListen()

	consumer, err := rmqtools.NewConsumer(app.rabbit, AmqpExchange, AmqpQueue)
	if err != nil {
		log.Panic(err)
	}
	err = consumer.Listen([]string{SignInKey, SignUpKey}, handlePayload)
	if err != nil {
		log.Panic(err)
	}
}

type AuthServer struct {
	auth.UnimplementedAuthServer
}

func (as *AuthServer) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	payload := req.GetPayload()
	log.Printf("sign in: incoming payload - %v\n\n", payload)

	// return response
	res := &auth.SignInResponse{Message: fmt.Sprintf("sign-in: get your payload - %v", payload)}
	return res, nil
}

func (as *AuthServer) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	payload := req.GetPayload()
	log.Printf("sign up: incoming payload - %v\n\n", payload)
	// return response
	res := &auth.SignUpResponse{Message: fmt.Sprintf("sign-up: get your payload - %v", payload)}
	return res, nil
}

func (app *App) gRPCListen() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("failed to listen tcp port - %s. Err - %s", gRPCPort, err.Error())
	}

	gRPCServer := grpc.NewServer()
	auth.RegisterAuthServer(gRPCServer, &AuthServer{})

	log.Printf("gRPC Server started on port %s", gRPCPort)

	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
