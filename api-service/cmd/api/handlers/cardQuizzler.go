package handlers

import (
	quizService "github.com/Salladin95/card-quizzler-microservices/api-service/card-quizzler"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
)

func (ah *apiHandlers) ProcessQuizResult(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "ProcessQuizResult")

	var dto entities.QuizResultDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	id := c.Param("moduleID")

	terms, err := lib.MarshalData(&dto.Terms)

	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		ProcessQuizResult(ctx, &quizService.ProcessQuizRequest{
			Terms:    terms,
			ModuleID: id,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("ProcessQuizResult", err)
	}
	return handleGRPCResponseNoContent(c, response)
}

func (ah *apiHandlers) CreateFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "CreateFolder")

	var dto entities.CreateFolderDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateFolder(ctx, &quizService.CreateFolderRequest{
			Payload: &quizService.CreateFolderPayload{
				Title:  dto.Title,
				UserID: dto.UserID,
			},
		})
	if err != nil {
		return goErrorHandler.OperationFailure("CreateFolder", err)
	}
	var unmarshalTo []entities.Folder
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) UpdateFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "UpdateFolder")

	id := c.Param("id")

	var dto entities.UpdateFolderDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		UpdateFolder(ctx, &quizService.UpdateFolderRequest{
			Payload: &quizService.UpdateFolderPayload{
				Title:    dto.Title,
				FolderID: id,
			},
		})
	if err != nil {
		return goErrorHandler.OperationFailure("UpdateFolder", err)
	}
	var unmarshalTo []entities.Folder
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) AddFolderToUser(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "AddFolderToUser")

	folderID := c.QueryParam("folderID")
	userID := c.QueryParam("userID")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		AddFolderToUser(ctx, &quizService.AddFolderToUserRequest{
			UserID:   userID,
			FolderID: folderID,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("AddFolderToUser", err)
	}
	var unmarshalTo entities.Folder
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) GetUserFolders(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetUserFolders")

	uid := c.Param("uid")

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetUserFolders(ctx, &quizService.RequestWithID{
			Id: uid,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUserFolders", err)
	}
	var unmarshalTo []entities.Folder
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) GetFolderByID(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetFolderByID")

	id := c.Param("id")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetFolderByID(ctx, &quizService.RequestWithID{
			Id: id,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("GetFolderByID", err)
	}
	var unmarshalTo entities.Folder
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) DeleteFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "DeleteFolder")

	id := c.Param("id")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		DeleteFolder(ctx, &quizService.RequestWithID{
			Id: id,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("DeleteFolder", err)
	}
	return handleGRPCResponseNoContent(c, response)
}

func (ah *apiHandlers) DeleteModuleFromFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "DeleteModuleFromFolder")

	// Retrieve the folderID and moduleID query parameters
	folderID := c.QueryParam("folderID")
	moduleID := c.QueryParam("moduleID")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		DeleteModuleFromFolder(ctx, &quizService.DeleteModuleFromFolderRequest{
			ModuleID: moduleID,
			FolderID: folderID,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("DeleteModuleFromFolder", err)
	}
	return handleGRPCResponseNoContent(c, response)
}

func (ah *apiHandlers) CreateModule(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "CreateModule")

	var dto entities.CreateModuleDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	terms, err := lib.MarshalData(&dto.Terms)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateModule(ctx, &quizService.CreateModuleRequest{
			Payload: &quizService.CreateModulePayload{
				Title:  dto.Title,
				UserID: dto.UserID,
				Terms:  terms,
			},
		})
	if err != nil {
		return goErrorHandler.OperationFailure("CreateModule", err)
	}
	var unmarshalTo entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) CreateModuleInFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "CreateModuleInFolder")

	folderID := c.Param("id")

	var dto entities.CreateModuleDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	terms, err := lib.MarshalData(&dto.Terms)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateModuleInFolder(ctx, &quizService.CreateModuleInFolderRequest{
			Payload: &quizService.CreateModulePayload{
				Title:  dto.Title,
				UserID: dto.UserID,
				Terms:  terms,
			},
			FolderID: folderID,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("CreateModuleInFolder", err)
	}
	var unmarshalTo entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) UpdateModule(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "UpdateModule")

	var dto entities.UpdateModuleDto
	if err := lib.BindBody(c, &dto); err != nil {
		return err
	}

	id := c.Param("id")

	newTerms, err := lib.MarshalData(&dto.NewTerms)
	if err != nil {
		return err
	}

	updatedTerms, err := lib.MarshalData(&dto.UpdatedTerms)
	if err != nil {
		return err
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		UpdateModule(ctx, &quizService.UpdateModuleRequest{
			Payload: &quizService.UpdateModulePayload{
				Title:        dto.Title,
				UpdatedTerms: updatedTerms,
				NewTerms:     newTerms,
				Id:           id,
			},
		})
	if err != nil {
		return goErrorHandler.OperationFailure("UpdateModule", err)
	}
	var unmarshalTo entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) GetUserModules(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetUserModules")

	uid := c.Param("uid")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetUserModules(ctx, &quizService.RequestWithID{
			Id: uid,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("GetUserModules", err)
	}
	var unmarshalTo []entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) GetModuleByID(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetModuleByID")

	id := c.Param("id")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetModuleByID(ctx, &quizService.RequestWithID{
			Id: id,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("GetModuleByID", err)
	}
	var unmarshalTo entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) AddModuleToUser(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "AddModuleToUser")

	moduleID := c.QueryParam("moduleID")
	userID := c.QueryParam("userID")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		AddModuleToUser(ctx, &quizService.AddModuleToUserRequest{
			UserID:   userID,
			ModuleID: moduleID,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("AddModuleToUser", err)
	}
	return handleGRPCResponseNoContent(c, response)
}

func (ah *apiHandlers) DeleteModule(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "DeleteModule")

	id := c.Param("id")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		DeleteModule(ctx, &quizService.RequestWithID{
			Id: id,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("DeleteModule", err)
	}
	return handleGRPCResponseNoContent(c, response)
}