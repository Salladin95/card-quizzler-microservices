package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"net/http"
)

// GetUserById retrieves a user based on the provided user ID.
func (us *UserServer) GetUserById(ctx context.Context, req *userService.ID) (*userService.Response, error) {
	// Log a message indicating the start of processing the GetUserById request
	us.log(ctx, "start processing GetUserById request", "info", "GetUserById")

	// Extract the user ID from the request
	id := req.GetId()

	// Fetch the user from the repository by ID
	fetchedUser, err := us.CachedRepo.GetById(ctx, id)
	if err != nil {
		// If an error occurs during user retrieval, build and return a failed response
		return buildFailedResponse(err)
	}

	// If user retrieval is successful, build and return a successful response
	return buildSuccessfulResponse(fetchedUser.ToResponse(), http.StatusOK, "user has been successfully fetched")
}

// GetUserByEmail retrieves a user based on the provided email.
func (us *UserServer) GetUserByEmail(ctx context.Context, req *userService.Email) (*userService.Response, error) {
	// Print a message indicating the start of processing the GetUserByEmail request
	us.log(ctx, "start processing GetUserByEmail request", "info", "GetUserByEmail")

	// Extract the email from the request
	email := req.GetEmail()

	// Fetch the user from the repository by email
	fetchedUser, err := us.CachedRepo.GetByEmail(ctx, email)
	if err != nil {
		// If an error occurs during user retrieval, build and return a failed response
		return buildFailedResponse(err)
	}

	// If user retrieval is successful, build and return a successful response
	return buildSuccessfulResponse(fetchedUser.ToResponse(), http.StatusOK, "user has been successfully fetched")
}

// UpdateEmail updates a user based on the provided UpdateEmailRequest.
func (us *UserServer) UpdateEmail(
	ctx context.Context,
	req *userService.UpdateEmailRequest,
) (*userService.Response, error) {
	// Print a message indicating the start of processing the update user request
	us.log(ctx, "start processing UpdateEmail request", "info", "UpdateEmail")

	// Extract payload from the gRPC request
	reqPayload := req.GetPayload()

	// Create a UpdateEmailDto from the request payload
	updateDto := user.UpdateEmailDto{
		Email: reqPayload.Email,
		Code:  reqPayload.Code,
	}

	// Verify the UpdateEmailDto structure
	err := updateDto.Verify()
	if err != nil {
		// Return a response with the mapped error status code and message if user update fails
		return buildFailedResponse(err)
	}

	existingUser, err := us.CachedRepo.GetById(ctx, reqPayload.Id)

	if err != nil {
		// Return a response with the mapped error status code and message if user update fails
		return buildFailedResponse(err)
	}

	// extract generated code
	cachedEmailCode, err := us.CachedRepo.GetEmailVerificationCode(ctx, existingUser.Email)

	if err != nil {
		// Return a response with the mapped error status code and message if user update fails
		return &userService.Response{Code: http.StatusBadRequest, Message: "You have to request verification firstly"}, nil
	}

	if cachedEmailCode.Code != int(updateDto.Code) {
		return &userService.Response{Code: 400, Message: "cached code and provided code don't match"}, nil
	}

	// Update the user by calling the UpdateEmail method in the repository
	updatedUser, err := us.CachedRepo.UpdateUser(
		ctx,
		reqPayload.Id,
		user.UpdateUserDto{Email: updateDto.Email},
	)
	if err != nil {
		// Return a response with the mapped error status code and message if user update fails
		return buildFailedResponse(err)
	}

	// Build and return a user response with a success code and message
	return buildSuccessfulResponse(updatedUser.ToResponse(), http.StatusCreated, "user has been updated")
}

// DeleteUser deletes a user based on the provided user ID.
func (us *UserServer) DeleteUser(ctx context.Context, req *userService.ID) (*userService.Response, error) {
	// Print a message indicating the start of processing the DeleteUser request.
	us.log(ctx, "start processing DeleteUser request", "info", "DeleteUser")

	// Extract the user ID from the request.
	id := req.GetId()

	// Call the repository's DeleteUser method to delete the user with the specified ID.
	err := us.CachedRepo.DeleteUser(ctx, id)
	if err != nil {
		// If an error occurs during the deletion process, build and return a failed response.
		return buildFailedResponse(err)
	}

	// If the user deletion is successful, return a successful response.
	return &userService.Response{Code: http.StatusNoContent, Message: "user has been deleted"}, nil
}
