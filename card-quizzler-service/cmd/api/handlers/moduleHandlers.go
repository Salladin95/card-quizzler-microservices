package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"net/http"
)

func (cq *CardQuizzlerServer) ProcessQuizResult(ctx context.Context, req *quizService.ProcessQuizRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "ProcessQuizResult")
	payload := req.GetTerms()
	mID := req.GetModuleID()
	moduleID, err := uuid.Parse(mID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	var dto entities.QuizResultDto
	if err := lib.UnmarshalData(payload, &dto); err != nil {
		return buildFailedResponse(err)
	}

	// Create a map to store terms and their answers
	answerMap := make(map[string]bool)
	for _, term := range dto.Terms {
		answerMap[term.ID] = term.Answer
	}

	module, err := cq.Repo.GetModuleByID(ctx, moduleID)

	// Update streaks and difficulty for each term in the module
	for _, term := range module.Terms {
		answer := answerMap[term.ID.String()]
		term.UpdateStreaksAndUpdateDifficulty(answer)
	}

	// Update the module with the updated terms
	if _, err := cq.Repo.UpdateModule(ctx, moduleID, entities.UpdateModuleDto{
		UpdatedTerms: module.Terms,
	}); err != nil {
		// Return a failed response if module update fails
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.MutateModuleKey, module)

	return buildSuccessfulResponse(nil, http.StatusOK, "Quiz has been processed")
}

func (cq *CardQuizzlerServer) AddModuleToUser(ctx context.Context, req *quizService.AddModuleToUserRequest) (*quizService.Response, error) {
	// Log the start of processing the gRPC request
	cq.log(ctx, "start processing grpc request", "info", "AddModuleToUser")

	// Extract module and user IDs from the request
	mID := req.GetModuleID()
	userID := req.GetUserID()

	if userID == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "User id is required"}, nil
	}

	moduleID, err := uuid.Parse(mID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Add module to user in the repository
	if err := cq.Repo.AddModuleToUser(userID, moduleID); err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(http.StatusNoContent, http.StatusOK, "Module is added to user")
}

func (cq *CardQuizzlerServer) CreateModule(ctx context.Context, req *quizService.CreateModuleRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "CreateModule")
	payload := req.GetPayload()

	// Unmarshal new terms from the payload
	var newTerms []entities.CreateTermDto
	if err := lib.UnmarshalData(payload.Terms, &newTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create a CreateModuleDto with extracted data
	createModuleDto := entities.CreateModuleDto{Title: payload.Title, UserID: payload.UserID, Terms: newTerms}

	if err := createModuleDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Create the module in the repository
	createdModule, err := cq.Repo.CreateModule(ctx, createModuleDto)
	if err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.CreateModuleKey, createdModule)

	return buildSuccessfulResponse(createdModule, http.StatusCreated, "module is created")
}

func (cq *CardQuizzlerServer) CreateModuleInFolder(ctx context.Context, req *quizService.CreateModuleInFolderRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "CreateModuleInFolder")
	payload := req.GetPayload()

	fID := req.GetFolderID()

	folderID, err := uuid.Parse(fID)

	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrBadRequest, err)
	}

	// Unmarshal new terms from the payload
	var newTerms []entities.CreateTermDto
	if err := lib.UnmarshalData(payload.Terms, &newTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create a CreateModuleDto with extracted data
	createModuleDto := entities.CreateModuleDto{Title: payload.Title, UserID: payload.UserID, Terms: newTerms}

	if err := createModuleDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Create the module in the repository
	createdModule, err := cq.Repo.CreateModule(ctx, createModuleDto)

	if err != nil {
		return buildFailedResponse(err)
	}

	if err := cq.Repo.AddModuleToFolder(ctx, folderID, createdModule.ID); err != nil {
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.CreateModuleKey, createdModule)
	cq.Broker.PushToQueue(ctx, constants.MutateFolderKey, models.Folder{ID: folderID})
	return buildSuccessfulResponse(createdModule, http.StatusCreated, "module is created")
}

func (cq *CardQuizzlerServer) UpdateModule(ctx context.Context, req *quizService.UpdateModuleRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "UpdateModule")

	payload := req.GetPayload()

	// Parse module ID from the payload
	moduleID, err := uuid.Parse(payload.Id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Unmarshal new terms from the payload
	var newTerms []entities.CreateTermDto
	if err := lib.UnmarshalData(payload.NewTerms, &newTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Unmarshal updated terms from the payload
	var updatedTerms []models.Term
	if err := lib.UnmarshalData(payload.UpdatedTerms, &updatedTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create an UpdateModuleDto with extracted data
	updateModuleDTO := entities.UpdateModuleDto{Title: payload.Title, NewTerms: newTerms, UpdatedTerms: updatedTerms}

	if err := updateModuleDTO.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Update the module in the repository
	updatedModule, err := cq.Repo.UpdateModule(ctx, moduleID, updateModuleDTO)
	if err != nil {
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.MutateModuleKey, updatedModule)
	return buildSuccessfulResponse(updatedModule, http.StatusOK, "module is updated")
}

func (cq *CardQuizzlerServer) DeleteModule(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "DeleteModule")

	id := req.GetId()
	moduleID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Delete the module from the repository
	if err := cq.Repo.DeleteModule(ctx, moduleID); err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.DeleteModuleKey, models.Module{ID: moduleID})
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Module is deleted")
}

func (cq *CardQuizzlerServer) GetUserModules(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "GetUserModules")

	uid := req.GetId()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}

	// Retrieve modules associated with the user from the repository
	modules, err := cq.Repo.GetModulesByUID(ctx, uid)
	if err != nil {
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.FetchUserModulesKey, modules)
	return buildSuccessfulResponse(modules, http.StatusOK, "requested modules")
}

func (cq *CardQuizzlerServer) GetModuleByID(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "GetModuleByID")

	id := req.GetId()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Retrieve the module by its ID from the repository
	module, err := cq.Repo.GetModuleByID(ctx, parsedID)
	if err != nil {
		// Return a failed response if retrieving the module fails
	}

	cq.Broker.PushToQueue(ctx, constants.FetchModuleKey, module)
	return buildSuccessfulResponse(module, http.StatusOK, "requested module")
}
