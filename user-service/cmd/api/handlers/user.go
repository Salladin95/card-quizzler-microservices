package handlers

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	model "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"net/http"
)

// GetUsers retrieves all users.
func (us *UserServer) GetUsers(ctx context.Context, _ *userService.EmptyRequest) (*userService.Response, error) {
	// Print a message indicating the start of processing the GetUsers request
	fmt.Println("******* user-service - start processing GetUsers request ********")

	// Fetch all users from the repository
	fetchedUsers, err := us.Repo.GetUsers(ctx)
	if err != nil {
		// If an error occurs during user retrieval, build and return a failed response
		return buildFailedResponse(err)
	}

	// If user retrieval is successful, build and return a successful response
	userResponses := model.ToUserResponses(fetchedUsers)
	return buildSuccessfulResponse(userResponses, http.StatusOK, "users have been successfully fetched")
}

// GetUserById retrieves a user based on the provided user ID.
func (us *UserServer) GetUserById(ctx context.Context, req *userService.ID) (*userService.Response, error) {
	// Print a message indicating the start of processing the GetUserById request
	fmt.Println("******* user-service - start processing GetUserById request ********")

	// Extract the user ID from the request
	id := req.GetId()

	// Fetch the user from the repository by ID
	fetchedUser, err := us.Repo.GetById(ctx, id)
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
	fmt.Println("******* user-service - start processing GetUserByEmail request ********")

	// Extract the email from the request
	email := req.GetEmail()

	// Fetch the user from the repository by email
	fetchedUser, err := us.Repo.GetByEmail(ctx, email)
	if err != nil {
		// If an error occurs during user retrieval, build and return a failed response
		return buildFailedResponse(err)
	}

	// If user retrieval is successful, build and return a successful response
	return buildSuccessfulResponse(fetchedUser.ToResponse(), http.StatusOK, "user has been successfully fetched")
}

// UpdateUser updates a user based on the provided UpdateUserRequest.
func (us *UserServer) UpdateUser(ctx context.Context, req *userService.UpdateUserRequest) (*userService.Response, error) {
	// Print a message indicating the start of processing the update user request
	fmt.Println("******* user-service - start processing update user request ********")

	// Extract payload from the gRPC request
	reqPayload := req.GetPayload()

	// Create a UpdateDto from the request payload
	updateDto := user.UpdateDto{
		Email:    reqPayload.Email,
		Password: reqPayload.Password,
		Name:     reqPayload.Name,
	}

	// Verify the UpdateDto structure
	err := updateDto.Verify()
	if err != nil {
		// Return a response with the mapped error status code and message if verification fails
		return &userService.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}

	// Update the user by calling the UpdateUser method in the repository
	updatedUser, err := us.Repo.UpdateUser(ctx, reqPayload.Id, updateDto)
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
	fmt.Println("******* user-service - start processing DeleteUser request ********")

	// Extract the user ID from the request.
	id := req.GetId()

	// Call the repository's DeleteUser method to delete the user with the specified ID.
	err := us.Repo.DeleteUser(ctx, id)
	if err != nil {
		// If an error occurs during the deletion process, build and return a failed response.
		return buildFailedResponse(err)
	}

	// If the user deletion is successful, return a successful response.
	return &userService.Response{Code: http.StatusNoContent, Message: "user has been deleted"}, nil
}
