package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/models"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/google/uuid"
	"net/http"
)

func (cq *CardQuizzlerServer) CreateFolder(ctx context.Context, req *quizService.CreateFolderRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "CreateFolderRequest")
	payload := req.GetPayload()
	createFolderDto := entities.CreateFolderDto{Title: payload.Title, UserID: payload.UserID}
	if err := createFolderDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}
	// Create the folder using the repository
	createdFolder, err := cq.Repo.CreateFolder(ctx, createFolderDto)
	if err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.CreateFolderKey, createdFolder)
	return buildSuccessfulResponse(createdFolder, http.StatusCreated, "created folder")
}

func (cq *CardQuizzlerServer) UpdateFolder(ctx context.Context, req *quizService.UpdateFolderRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "UpdateFolder")
	payload := req.GetPayload()
	folderID, err := uuid.Parse(payload.FolderID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	updateFolderDto := entities.UpdateFolderDto{Title: payload.Title}
	if err := updateFolderDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}
	updateFolder, err := cq.Repo.UpdateFolder(ctx, folderID, updateFolderDto)
	if err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.MutateFolderKey, updateFolder)
	return buildSuccessfulResponse(updateFolder, http.StatusOK, "updated folder")
}

func (cq *CardQuizzlerServer) AddFolderToUser(ctx context.Context, req *quizService.AddFolderToUserRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "AddFolderToUser")
	fID := req.GetFolderID()
	userID := req.GetUserID()
	folderID, err := uuid.Parse(fID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	if userID == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "User id is required"}, nil
	}
	if err := cq.Repo.AddFolderToUser(userID, folderID); err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(http.StatusNoContent, http.StatusOK, "folder is added to user")
}

func (cq *CardQuizzlerServer) DeleteFolder(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "DeleteFolder")
	id := req.GetId()
	folderID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	if err := cq.Repo.DeleteFolder(ctx, folderID); err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.DeleteFolderKey, models.Folder{ID: folderID})
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Folder is deleted")
}

func (cq *CardQuizzlerServer) DeleteModuleFromFolder(ctx context.Context, req *quizService.DeleteModuleFromFolderRequest) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "DeleteModuleFromFolder")
	fID := req.GetFolderID()
	mID := req.GetModuleID()
	folderID, err := uuid.Parse(fID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	moduleID, err := uuid.Parse(mID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	if err := cq.Repo.DeleteModuleFromFolder(ctx, folderID, moduleID); err != nil {
		return buildFailedResponse(err)
	}
	cq.Broker.PushToQueue(ctx, constants.MutateFolderAndModule, entities.FolderAndModuleIDS{
		FolderID: folderID,
		ModuleID: moduleID,
	})
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Module is deleted from folder ")
}

func (cq *CardQuizzlerServer) GetUserFolders(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "GetUserFolders")
	uid := req.GetId()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}

	folders, err := cq.Repo.GetFoldersByUID(ctx, uid)
	if err != nil {
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.FetchUserFoldersKey, folders)
	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetUserFolderByID(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	cq.log(ctx, "start processing grpc request", "info", "GetUserFolderByID")
	id := req.GetId()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	folder, err := cq.Repo.GetFolderByID(ctx, parsedID)
	if err != nil {
		return buildFailedResponse(err)
	}

	cq.Broker.PushToQueue(ctx, constants.FetchFolderKey, folder)
	return buildSuccessfulResponse(folder, http.StatusOK, "requested folder")
}