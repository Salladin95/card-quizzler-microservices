package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/cmd/api/repositories"
	quizService "github.com/Salladin95/card-quizzler-microservices/card-quizzler-service/proto"
	"github.com/google/uuid"
	"net/http"
)

func (cq *CardQuizzlerServer) CreateFolder(ctx context.Context, req *quizService.CreateFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[CreateFolder] Start processing grpc request")
	payload := req.GetPayload()
	createFolderDto := entities.CreateFolderDto{IsOpen: payload.IsOpen, Title: payload.Title, UserID: payload.UserID}
	if err := createFolderDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}
	// Create the folder using the repository
	createdFolder, err := cq.Repo.CreateFolder(ctx, createFolderDto)
	if err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(createdFolder, http.StatusCreated, "created folder")
}

func (cq *CardQuizzlerServer) UpdateFolder(ctx context.Context, req *quizService.UpdateFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[UpdateFolder] Start processing grpc request")
	payload := req.GetPayload()
	folderID, err := uuid.Parse(payload.FolderID)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	updateFolderDto := entities.UpdateFolderDto{Title: payload.Title, IsOpen: payload.IsOpen}
	if err := updateFolderDto.Verify(); err != nil {
		return buildFailedResponse(err)
	}
	updateFolder, err := cq.Repo.UpdateFolder(repositories.UpdateFolderPayload{
		Ctx:      ctx,
		FolderID: folderID,
		Dto:      updateFolderDto,
	})
	if err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(updateFolder, http.StatusOK, "updated folder")
}

func (cq *CardQuizzlerServer) AddFolderToUser(ctx context.Context, req *quizService.AddFolderToUserRequest) (*quizService.Response, error) {
	lib.LogInfo("[AddFolderToUser] Start processing grpc request")
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
	lib.LogInfo("[DeleteFolder] Start processing grpc request")
	id := req.GetId()
	folderID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}
	if err := cq.Repo.DeleteFolder(ctx, folderID); err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Folder is deleted")
}

func (cq *CardQuizzlerServer) DeleteModuleFromFolder(ctx context.Context, req *quizService.DeleteModuleFromFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[DeleteModuleFromFolder] Start processing grpc request")
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
	if err := cq.Repo.DeleteModuleFromFolder(repositories.FolderModuleAssociation{
		Ctx:      ctx,
		FolderID: folderID,
		ModuleID: moduleID,
	}); err != nil {
		return buildFailedResponse(err)
	}
	return buildSuccessfulResponse(nil, http.StatusNoContent, "Module is deleted from folder ")
}

func (cq *CardQuizzlerServer) GetUserFolders(ctx context.Context, req *quizService.GetUserFoldersRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetUserFolders] Start processing grpc request")
	uid := req.GetId()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}
	payload := req.GetPayload()

	folders, err := cq.Repo.GetFoldersByUID(repositories.UidSortPayload{
		Ctx:    ctx,
		Uid:    uid,
		Limit:  payload.Limit,
		Page:   payload.Page,
		SortBy: payload.SortBy,
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetOpenFolders(ctx context.Context, req *quizService.GetOpenFoldersRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetOpenFolders] Start processing grpc request")
	payload := req.GetPayload()

	folders, err := cq.Repo.GetOpenFolders(repositories.UidSortPayload{
		Ctx:    ctx,
		Limit:  payload.Limit,
		Page:   payload.Page,
		SortBy: payload.SortBy,
	})
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetFolderByID(ctx context.Context, req *quizService.RequestWithID) (*quizService.Response, error) {
	lib.LogInfo("[GetFolderByID] Start processing grpc request")
	id := req.GetId()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	folder, err := cq.Repo.GetFolderByID(ctx, parsedID)
	if err != nil {
		return buildFailedResponse(err)
	}

	return buildSuccessfulResponse(folder, http.StatusOK, "requested folder")
}
