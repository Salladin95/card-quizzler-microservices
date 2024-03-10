package handlers

import (
	"errors"
	"fmt"
	quizService "github.com/Salladin95/card-quizzler-microservices/api-service/card-quizzler"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/cacheManager"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	"github.com/Salladin95/goErrorHandler"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ah *apiHandlers) GetUserFolders(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetUserFolders")

	limit := ParseInt(c.QueryParam("limit"), foldersDefaultLimit)
	page := ParseInt(c.QueryParam("page"), 1)
	sortBy := ParseSortBy(
		c.QueryParam("sortBy"),
		"asc",
		"created_at",
		FolderKeysMap,
	)

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	var folders []entities.Folder
	err := ah.cacheManager.ReadCacheByKeys(
		&folders,
		cacheManager.FoldersKey(uid),
		fmt.Sprintf("%d:%d:%s", limit, page, sortBy),
	)

	if err == nil {
		ah.log(ctx, "retrieved from cache", "info", "GetUserFolders")
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested folders", Data: folders})
	}

	// Obtain a gRPC client connection using the GetGRPCClientConn method from apiHandlers.
	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	// Make a gRPC call to the SignIn method of the Auth service
	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetUserFolders(ctx, &quizService.GetUserFoldersRequest{
			Id: uid,
			Payload: &quizService.SortOptions{
				Limit:  limit,
				Page:   page,
				SortBy: sortBy,
			},
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

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}

	var folder entities.Folder
	err := ah.cacheManager.ReadCacheByKeys(
		&folder,
		cacheManager.FolderKey(claims.Id),
		cacheManager.FolderKey(id),
	)

	if err == nil {
		ah.log(ctx, "retrieved from cache", "info", "GetFolderByID")
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested folder", Data: folder})
	}

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
	return handleGRPCResponse(c, response, &folder)
}

func (ah *apiHandlers) GetUserModules(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetUserModules")

	limit := ParseInt(c.QueryParam("limit"), modulesDefaultLimit)
	page := ParseInt(c.QueryParam("page"), 1)
	sortBy := ParseSortBy(
		c.QueryParam("sortBy"),
		"asc",
		"created_at",
		ModuleKeysMap,
	)

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	var modules []entities.Module
	err := ah.cacheManager.ReadCacheByKeys(
		&modules,
		cacheManager.ModulesKey(uid),
		fmt.Sprintf("%d:%d:%s", limit, page, sortBy),
	)

	if err == nil {
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested modules", Data: modules})
	}

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetUserModules(ctx, &quizService.GetUserModulesRequest{
			Id: uid,
			Payload: &quizService.SortOptions{
				Limit:  limit,
				Page:   page,
				SortBy: sortBy,
			},
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

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}

	var module entities.Module
	err := ah.cacheManager.ReadCacheByKeys(
		&module,
		cacheManager.ModuleKey(claims.Id),
		cacheManager.ModuleKey(id),
	)

	if err == nil {
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested module", Data: module})
	}

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

func (ah *apiHandlers) GetDifficultModules(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "GetDifficultModules")

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	var modules []entities.Module
	err := ah.cacheManager.ReadCacheByKeys(
		&modules,
		cacheManager.ModulesKey(claims.Id),
		cacheManager.DifficultModules,
	)

	if err == nil {
		return c.JSON(http.StatusOK, entities.JsonResponse{Message: "Requested modules", Data: modules})
	}

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		GetDifficultModulesByUID(ctx, &quizService.GetDifficultModulesRequest{
			Id: uid,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("GetDifficultModules", err)
	}
	var unmarshalTo []entities.Module
	return handleGRPCResponse(c, response, unmarshalTo)
}

func (ah *apiHandlers) ProcessQuizResult(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "ProcessQuizResult")

	var dto entities.QuizResultDto
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
		ProcessQuizResult(ctx, &quizService.ProcessQuizRequest{
			Terms: terms,
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

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateFolder(ctx, &quizService.CreateFolderRequest{
			Payload: &quizService.CreateFolderPayload{
				Title:  dto.Title,
				UserID: uid,
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

func (ah *apiHandlers) AddModuleToFolder(c echo.Context) error {
	ctx := c.Request().Context()
	ah.log(ctx, "start processing request", "info", "AddModuleToFolder")

	folderID := c.QueryParam("folderID")
	moduleID := c.QueryParam("moduleID")

	clientConn, err := ah.GetGRPCClientConn(ah.config.AppCfg.CardQuizServiceUrl)
	defer clientConn.Close() // Ensure the gRPC client connection is closed when done.
	if err != nil {
		return err // Return an error if obtaining the client connection fails.
	}

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		AddModuleToFolder(ctx, &quizService.AddModuleToFolderRequest{
			ModuleID: moduleID,
			FolderID: folderID,
		})
	if err != nil {
		return goErrorHandler.OperationFailure("AddModuleToFolder", err)
	}

	return handleGRPCResponseNoContent(c, response)
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

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateModule(ctx, &quizService.CreateModuleRequest{
			Payload: &quizService.CreateModulePayload{
				Title:  dto.Title,
				UserID: uid,
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

	// Retrieve user claims from the context
	claims, ok := c.Get("user").(*lib.JwtUserClaims)
	if !ok {
		return goErrorHandler.NewError(
			goErrorHandler.ErrUnauthorized,
			errors.New("failed to cast claims"),
		)
	}
	uid := claims.Id

	response, err := quizService.
		NewCardQuizzlerServiceClient(clientConn).
		CreateModuleInFolder(ctx, &quizService.CreateModuleInFolderRequest{
			Payload: &quizService.CreateModulePayload{
				Title:  dto.Title,
				UserID: uid,
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
