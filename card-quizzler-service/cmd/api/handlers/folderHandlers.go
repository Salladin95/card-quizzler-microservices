package handlers

import (
	"context"
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

func (cq *CardQuizzlerServer) CreateFolder(ctx context.Context, req *quizService.CreateFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[CreateFolder] Start processing grpc request")
	payload := req.GetPayload()
	secureAccess := payload.SecureAccess

	createFolderDto := entities.CreateFolderDto{
		SecureAccess: entities.SecureAccess{
			Password: secureAccess.Password,
			Access:   models.AccessType(strings.ToLower(secureAccess.Access)),
		},
		Title:  payload.Title,
		UserID: payload.UserID,
	}

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
	uid := payload.GetUid()

	folderID, err := lib.ParseUUID(payload.FolderID)
	if err != nil {
		return buildFailedResponse(err)
	}

	folder, err := cq.Repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return buildFailedResponse(err)
	}

	if err := checkOwnership(uid, folder.UserID); err != nil {
		return buildFailedResponse(err)
	}

	secureAccess := entities.SecureAccess{
		Access:   models.AccessType(strings.ToLower(payload.SecureAccess.Access)),
		Password: payload.SecureAccess.Password,
	}

	if err := CheckPassword(
		secureAccess,
		folder.Access,
	); err != nil {
		return buildFailedResponse(err)
	}

	updateFolderDto := entities.UpdateFolderDto{
		Title:    payload.Title,
		Access:   secureAccess.Access,
		Password: secureAccess.Password,
	}

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
	password := req.GetPassword()
	folderID, err := uuid.Parse(fID)

	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	uid := req.GetUserID()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "User id is required"}, nil
	}

	folder, err := cq.Repo.GetFolderByID(ctx, folderID)
	if err != nil {
		return buildFailedResponse(err)
	}

	if folder.Access == models.AccessOnlyMe {
		return buildFailedResponse(goErrorHandler.ForbiddenError())
	}

	if folder.Access == models.AccessPassword {
		if err := checkPassword(&folder, password); err != nil {
			return buildFailedResponse(err)
		}
	}

	if err := cq.Repo.AddFolderToUser(uid, folderID); err != nil {
		return buildFailedResponse(err)
	}
	return buildNoContentResponse(http.StatusNoContent, "folder is added to user")
}

func (cq *CardQuizzlerServer) DeleteFolder(ctx context.Context, req *quizService.RequestWithIdAndUID) (*quizService.Response, error) {
	lib.LogInfo("[DeleteFolder] Start processing grpc request")
	id := req.GetId()
	uid := req.GetUid()
	folderID, err := uuid.Parse(id)

	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	if err := cq.checkFolderOwnership(ctx, uid, folderID); err != nil {
		return buildFailedResponse(err)
	}

	if err := cq.Repo.DeleteFolder(ctx, folderID); err != nil {
		return buildFailedResponse(err)
	}

	return buildNoContentResponse(http.StatusNoContent, "Folder is deleted")
}

func (cq *CardQuizzlerServer) DeleteModuleFromFolder(ctx context.Context, req *quizService.DeleteModuleFromFolderRequest) (*quizService.Response, error) {
	lib.LogInfo("[DeleteModuleFromFolder] Start processing grpc request")
	fID := req.GetFolderID()
	mID := req.GetModuleID()
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

	folders, err := cq.Repo.GetFoldersByUID(repositories.GetByUIDPayload{
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

	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetFoldersByTitle(ctx context.Context, req *quizService.GetByTitleRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetFoldersByTitle] Start processing grpc request")
	title := req.GetTitle()
	if title == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "title is missing"}, nil
	}
	uid := req.GetUid()
	if uid == "" {
		return &quizService.Response{Code: http.StatusBadRequest, Message: "user id is missing"}, nil
	}
	sortOptions := req.GetSortOptions()

	folders, err := cq.Repo.GetFoldersByTitle(repositories.GetByTitlePayload{
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

	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetOpenFolders(ctx context.Context, req *quizService.GetOpenFoldersRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetOpenFolders] Start processing grpc request")
	payload := req.GetPayload()

	folders, err := cq.Repo.GetOpenFolders(repositories.GetByUIDPayload{
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

	return buildSuccessfulResponse(folders, http.StatusOK, "requested folders")
}

func (cq *CardQuizzlerServer) GetFolderByID(ctx context.Context, req *quizService.GetByIDRequest) (*quizService.Response, error) {
	lib.LogInfo("[GetFolderByID] Start processing grpc request")
	id := req.GetId()
	uid := req.GetUid()
	password := req.GetPassword()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return &quizService.Response{Code: http.StatusBadRequest, Message: err.Error()}, nil
	}

	folder, err := cq.Repo.GetFolderByID(ctx, parsedID)
	if err != nil {
		return buildFailedResponse(err)
	}

	return handleAccessValidation(&folder, uid, password)
}
