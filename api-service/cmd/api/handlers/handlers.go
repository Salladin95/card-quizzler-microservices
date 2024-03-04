package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ApiHandlersInterface defines the interface for api-service-related HTTP handlers.
type ApiHandlersInterface interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	Refresh(c echo.Context) error
	UpdateEmail(c echo.Context) error
	UpdatePassword(c echo.Context) error
	ResetPassword(c echo.Context) error
	GetUserById(c echo.Context) error
	GetProfile(c echo.Context) error
	RequestEmailVerification(c echo.Context) error
	ProcessQuizResult(c echo.Context) error
	CreateFolder(c echo.Context) error
	UpdateFolder(c echo.Context) error
	AddFolderToUser(c echo.Context) error
	GetUserFolders(c echo.Context) error
	GetFolderByID(c echo.Context) error
	DeleteFolder(c echo.Context) error
	DeleteModuleFromFolder(c echo.Context) error
	CreateModule(c echo.Context) error
	CreateModuleInFolder(c echo.Context) error
	UpdateModule(c echo.Context) error
	GetUserModules(c echo.Context) error
	GetDifficultModules(c echo.Context) error
	GetModuleByID(c echo.Context) error
	AddModuleToUser(c echo.Context) error
	AddModuleToFolder(c echo.Context) error
	DeleteModule(c echo.Context) error
}

// apiHandlers implements the ApiHandlersInterface.
type apiHandlers struct {
	broker       rmqtools.MessageBroker
	config       *config.Config
	cacheManager cacheManager.CacheManager
}

// NewHandlers creates a new instance of ApiHandlersInterface.
func NewHandlers(
	cfg *config.Config,
	broker rmqtools.MessageBroker,
	cacheManager cacheManager.CacheManager,
) ApiHandlersInterface {
	return &apiHandlers{
		broker:       broker,
		config:       cfg,
		cacheManager: cacheManager,
	}
}

// GetGRPCClientConn establishes a gRPC client connection using the specified URL and returns the connection.
func (ah *apiHandlers) GetGRPCClientConn(url string) (*grpc.ClientConn, error) {
	// Dial a gRPC server using the provided URL and insecure transport credentials
	conn, err := grpc.Dial(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		// HandleEvent the error and return an OperationFailure error using the goErrorHandler package
		return nil, goErrorHandler.OperationFailure("connect to gRPC", err)
	}

	return conn, nil
}

// log sends a log Message to the Message broker.
func (ah *apiHandlers) log(ctx context.Context, message, level, method string) {
	var log entities.LogMessage

	// Push log Message to the Message broker
	ah.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log Message with provided details
		log.GenerateLog(message, level, method, "http handlers"),
	)
}
