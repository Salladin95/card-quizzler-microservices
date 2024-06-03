package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (cq *CardQuizzlerServer) ProcessQuizResult(ctx context.Context, req *quizService.ProcessQuizRequest) (*quizService.Response, error) {
	lib.LogInfo("[ProcessQuizResult] Start processing grpc request", "info", "ProcessQuizResult")
	payload := req.GetTerms()
	id := req.GetModuleID()

	// if there's id we want to update updated_at, so we can fetch recent opened modules
	if id != "" {
		moduleID, err := uuid.Parse(id)
		if err != nil {
			return buildFailedResponse(err)
		}

		if _, err = cq.Repo.UpdateModule(repositories.UpdateModulePayload{
			Ctx:      ctx,
			ModuleID: moduleID,
			Dto:      entities.UpdateModuleDto{},
		}); err != nil {
			return buildFailedResponse(err)
		}
	}

	var resultTerms []entities.ResultTerm
	if err := lib.UnmarshalData(payload, &resultTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create a map to store terms and their answers
	answersMap := make(map[string]bool)
	for _, term := range resultTerms {
		answersMap[term.ID] = term.Answer
	}

	var termIDS []uuid.UUID
	for _, term := range resultTerms {
		id, err := uuid.Parse(term.ID)
		if err != nil {
			return buildFailedResponse(goErrorHandler.OperationFailure("parse UUID", err))
		}
		termIDS = append(termIDS, id)
	}

	terms, err := cq.Repo.GetTerms(termIDS)
	if err != nil {
		return buildFailedResponse(err)
	}

	// Update streaks and difficulty for each term in the module
	for i := range terms {
		answer := answersMap[terms[i].ID.String()]
		models.UpdateStreaksAndUpdateDifficulty(&terms[i], answer)
	}

	// Update the module with the updated terms
	if err := cq.Repo.UpdateTerms(ctx, terms); err != nil {
		// Return a failed response if module update fails
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(nil, http.StatusOK, "Quiz has been processed")
}

func (cq *CardQuizzlerServer) AddModuleToUser(ctx context.Context, req *quizService.AddModuleToUserRequest) (*quizService.Response, error) {
	// Log the start of processing the gRPC request
	lib.LogInfo("[AddModuleToUser] Start processing grpc request", "info", "AddModuleToUser")

	// Extract module and user IDs from the request
	mID := req.GetModuleID()
	uid := req.GetUserID()
	password := req.GetPassword()

	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "User id is required"}, nil
	}

	moduleID, err := uuid.Parse(mID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	module, err := cq.Repo.GetModuleByID(ctx, moduleID)
	if err != nil {
		return buildFailedResponse(err)
	}

	if module.Access == models.AccessOnlyMe {
		return buildFailedResponse(goErrorHandler.ForbiddenError())
	}

	if module.Access == models.AccessPassword {
		if err := checkPassword(&module, password); err != nil {
			return buildFailedResponse(err)
		}
	}

	// Add module to user in the repository
	if err := cq.Repo.AddModuleToUser(uid, moduleID); err != nil {
		return buildFailedResponse(err)
	}

	return buildNoContentResponse(http.StatusNoContent, "Module is added to user")
}

func (cq *CardQuizzlerServer) AddModuleToFolder(ctx context.Context, req *quizService.AddModuleToFolderRequest) (*quizService.Response, error) {
	// Log the start of processing the gRPC request
	lib.LogInfo("[AddModuleToFolder] Start processing grpc request", "info", "AddModuleToUser")

	// Extract module and user IDs from the request
	mID := req.GetModuleID()
	fID := req.GetFolderID()
	uid := req.GetUid()

	folderID, err := uuid.Parse(fID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	moduleID, err := uuid.Parse(mID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	if err := cq.checkFolderOwnership(ctx, uid, folderID); err != nil {
		return buildFailedResponse(err)
	}

	if err := cq.checkModuleOwnership(ctx, uid, moduleID); err != nil {
		return buildFailedResponse(err)
	}

	// Add module to user in the repository
	if err := cq.Repo.AddModuleToFolder(repositories.FolderModuleAssociation{
		Ctx:      ctx,
		FolderID: folderID,
		ModuleID: moduleID,
	}); err != nil {
		return buildFailedResponse(err)
	}

	return buildNoContentResponse(http.StatusNoContent, "Module is added to folder")
}

func (cq *CardQuizzlerServer) CreateModule(ctx context.Context, req *quizService.CreateModuleRequest) (*quizService.Response, error) {
	lib.LogInfo("[CreateModule] Start processing grpc request", "info", "CreateModule")
	payload := req.GetPayload()
	secureAccess := payload.SecureAccess

	// Unmarshal new terms from the payload
	var newTerms []entities.CreateTermDto
	if err := lib.UnmarshalData(payload.Terms, &newTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create a CreateModuleDto with extracted data
	createModuleDto := entities.CreateModuleDto{
		Title: payload.Title,
		SecureAccess: entities.SecureAccess{
			Password: secureAccess.Password,
			Access:   models.AccessType(strings.ToLower(secureAccess.Access)),
		},
		UserID: payload.UserID,
		Terms:  newTerms,
	}

	if err := createModuleDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Create the module in the repository
	createdModule, err := cq.Repo.CreateModule(ctx, createModuleDto)
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(createdModule, http.StatusCreated, "module is created")
}

func (cq *CardQuizzlerServer) CreateModuleInFolder(ctx context.Context, req *quizService.CreateModuleInFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[CreateModuleInFolder] Start processing grpc request", "info", "CreateModuleInFolder")
	payload := req.GetPayload()

	fID := req.GetFolderID()
	uid := payload.UserID

	folderID, err := uuid.Parse(fID)

	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrBadRequest, err)
	}

	if err := cq.checkFolderOwnership(ctx, uid, folderID); err != nil {
		return buildFailedResponse(err)
	}

	// Unmarshal new terms from the payload
	var newTerms []entities.CreateTermDto
	if err := lib.UnmarshalData(payload.Terms, &newTerms); err != nil {
		return buildFailedResponse(err)
	}

	// Create a CreateModuleDto with extracted data
	createModuleDto := entities.CreateModuleDto{Title: payload.Title, UserID: uid, Terms: newTerms}

	if err := createModuleDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Create the module in the repository
	createdModule, err := cq.Repo.CreateModule(ctx, createModuleDto)

	if err != nil {
		return buildFailedResponse(err)
	}

	if err := cq.Repo.AddModuleToFolder(repositories.FolderModuleAssociation{
		Ctx:      ctx,
		FolderID: folderID,
		ModuleID: createdModule.ID,
	}); err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(createdModule, http.StatusCreated, "module is created")
}

func (cq *CardQuizzlerServer) UpdateModule(ctx context.Context, req *quizService.UpdateModuleRequest) (*quizService.Response, error) {
	lib.LogInfo("[UpdateModule] Start processing grpc request", "info", "UpdateModule")
	payload := req.GetPayload()
	uid := payload.GetUid()

	moduleID, err := lib.ParseUUID(payload.Id)
	if err != nil {
		return buildFailedResponse(err)
	}

	module, err := cq.Repo.GetModuleByID(ctx, moduleID)
	if err != nil {
		return buildFailedResponse(err)
	}

	if err := checkOwnership(uid, module.UserID); err != nil {
		return buildFailedResponse(err)
	}

	secureAccess := entities.SecureAccess{
		Access:   models.AccessType(strings.ToLower(payload.SecureAccess.Access)),
		Password: payload.SecureAccess.Password,
	}

	fmt.Println("<<<<<<<<<<<<<<<<<")

	if err := CheckPassword(
		secureAccess,
		module.Access,
	); err != nil {
		return buildFailedResponse(err)
	}

	fmt.Println(">>>>>>>>>>>>")

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

	updateModuleDto := entities.UpdateModuleDto{
		Title:        payload.Title,
		Access:       secureAccess.Access,
		Password:     secureAccess.Password,
		NewTerms:     newTerms,
		UpdatedTerms: updatedTerms,
	}

	if err := updateModuleDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Update the module in the repository
	updatedModule, err := cq.Repo.UpdateModule(repositories.UpdateModulePayload{
		Ctx:      ctx,
		ModuleID: moduleID,
		Dto:      updateModuleDto,
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(updatedModule, http.StatusOK, "module is updated")
}

func (cq *CardQuizzlerServer) DeleteModule(ctx context.Context, req *quizService.RequestWithIdAndUID) (*quizService.Response, error) {
	lib.LogInfo("[DeleteModule] Start processing grpc request", "info", "DeleteModule")

	id := req.GetId()
	uid := req.GetUid()
	moduleID, err := uuid.Parse(id)

	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	if err := cq.checkModuleOwnership(ctx, uid, moduleID); err != nil {
		return buildFailedResponse(err)
	}

	// Delete the module from the repository
	if err := cq.Repo.DeleteModule(ctx, moduleID); err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Module is deleted")
}

func (cq *CardQuizzlerServer) GetUserModules(ctx context.Context, req *quizService.GetUserModulesRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetUserModules] Start processing grpc request", "info", "GetUserModules")

	uid := req.GetId()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}

	payload := req.GetPayload()

	// Retrieve modules associated with the user from the repository
	modules, err := cq.Repo.GetModulesByUID(repositories.GetByUIDPayload{
		SortPayload: repositories.SortPayload{
			Ctx:    ctx,
			Uid:    uid,
			Limit:  payload.Limit,
			Page:   payload.Page,
			SortBy: payload.SortBy,
		},
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(modules, http.StatusOK, "requested modules")
}

func (cq *CardQuizzlerServer) GetModulesByTitle(ctx context.Context, req *quizService.GetByTitleRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetModulesByTitle] Start processing grpc request")
	title := req.GetTitle()
	if title == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "title is missing"}, nil
	}
	uid := req.GetUid()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}
	sortOptions := req.GetSortOptions()

	modules, err := cq.Repo.GetModulesByTitle(repositories.GetByTitlePayload{
		Title: title,
		SortPayload: repositories.SortPayload{
			Ctx:    ctx,
			Uid:    uid,
			Limit:  sortOptions.Limit,
			Page:   sortOptions.Page,
			SortBy: sortOptions.SortBy,
		},
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(modules, http.StatusOK, "requested modules")
}

func (cq *CardQuizzlerServer) GetOpenModules(ctx context.Context, req *quizService.GetOpenModulesRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetUserModules] Start processing grpc request", "info", "GetOpenModules")

	payload := req.GetPayload()

	// Retrieve modules associated with the user from the repository
	modules, err := cq.Repo.GetOpenModules(repositories.GetByUIDPayload{
		SortPayload: repositories.SortPayload{
			Ctx:    ctx,
			Limit:  payload.Limit,
			Page:   payload.Page,
			SortBy: payload.SortBy,
		},
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(modules, http.StatusOK, "requested modules")
}

func (cq *CardQuizzlerServer) GetDifficultModulesByUID(ctx context.Context, req *quizService.GetDifficultModulesRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetDifficultModulesByUID] Start processing grpc request", "info", "GetDifficultUserModules")

	uid := req.GetId()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}

	// Retrieve modules associated with the user from the repository
	modules, err := cq.Repo.GetDifficultModulesByUID(ctx, uid)
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(modules, http.StatusOK, "requested modules")
}

func (cq *CardQuizzlerServer) GetModuleByID(ctx context.Context, req *quizService.GetByIDRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetModuleByID] Start processing grpc request", "info", "GetModuleByID")

	uid := req.GetUid()
	password := req.GetPassword()

	id := req.GetId()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	// Retrieve the module by its ID from the repository
	module, err := cq.Repo.GetModuleByID(ctx, parsedID)
	if err != nil {
		return buildFailedResponse(err)
	}

	return handleAccessValidation(&module, uid, password)
}

func (cq *CardQuizzlerServer) UpdateTerm(ctx context.Context, req *quizService.UpdaterTermRequest) (*quizService.Response, error) {
	lib.LogInfo("[UpdateTerm] Start processing grpc request", "info", "UpdateTerm")
	termID, err := uuid.Parse(req.GetId())
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	moduleID, err := uuid.Parse(req.GetModuleID())
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	updateTermDto := entities.UpdateTermDto{
		Id:          termID,
		ModuleID:    moduleID,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}

	if err := updateTermDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}

	// Create the module in the repository
	if err := cq.Repo.UpdateTerm(ctx, updateTermDto); err != nil {
		return buildFailedResponse(err)
	}

	return &quizService.Response{
		Data:    nil,
		Code:    http.StatusOK,
		Message: "Term has been updated",
	}, nil
}
