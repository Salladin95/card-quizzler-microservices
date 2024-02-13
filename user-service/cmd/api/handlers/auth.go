package handlers

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/lib"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	userService "github.com/Salladin95/card-quizzler-microservices/user-service/proto"
	"net/http"
)

// SignIn is a gRPC service method that processes a sign-in request.
// It takes a context and a gRPC request (user-service.SignInRequest) as input.
// The function returns a gRPC response (user-service.Response) and an error.
// The response includes a status code and a message.
func (us *UserServer) SignIn(ctx context.Context, req *userService.SignInRequest) (*userService.Response, error) {
	// Print a message indicating the start of processing the sign-in request
	us.log(ctx, "start processing signIn request", "info", "signIn")
	// Extract payload from the gRPC request
	reqPayload := req.GetPayload()
	// Create a SignInDto from the request payload
	signInDto := user.SignInDto{Email: reqPayload.Email, Password: reqPayload.Password}
	// Verify the SignInDto structure
	err := signInDto.Verify()
	if err != nil {
		// Return a response with the mapped error status code and message if verification fails
		return &userService.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}

	// Check if the fetched user's email matches the sign-in email and verify the password
	fetchedUser, err := us.Repo.GetByEmail(ctx, signInDto.Email)

	// Fetch user by email from the repository
	if err != nil {
		// Return a response with the mapped error status code and message if fetching user fails
		return &userService.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	isPasswordInvalid := lib.CompareHashAndPassword(fetchedUser.Password, signInDto.Password)
	if fetchedUser.Email != signInDto.Email || isPasswordInvalid != nil {
		// Return a response with the mapped error status code and message if authentication fails
		return buildFailedResponse(err)
	}
	// Build and return a user response with a success code and message
	return buildSuccessfulResponse(fetchedUser.ToResponse(), http.StatusOK, "user has signed in")
}

// SignUp is a gRPC service method that handles the user sign-up request.
// It takes a context and a gRPC request (user-service.SignUpRequest) as input.
// The function returns a gRPC response (user-service.Response) and an error.
// The response includes a status code and a message.
func (us *UserServer) SignUp(ctx context.Context, req *userService.SignUpRequest) (*userService.Response, error) {
	// Print a message indicating the start of processing the sign-up request
	us.log(ctx, "start processing signUp request", "info", "signIn")
	// Extract payload from the gRPC request
	reqPayload := req.GetPayload()
	// Create a SignUpDto from the request payload
	signUpDto := user.SignUpDto{
		Email:    reqPayload.Email,
		Password: reqPayload.Password,
		Name:     reqPayload.Name,
		Birthday: reqPayload.Birthday,
	}
	// Verify the SignUpDto structure
	err := signUpDto.Verify()
	if err != nil {
		// Return a response with the mapped error status code and message if verification fails
		return &userService.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Create a new user by calling the CreateUser method in the repository
	newUser, err := us.Repo.CreateUser(ctx, signUpDto)
	if err != nil {
		// Return a response with the mapped error status code and message if user creation fails
		return buildFailedResponse(err)
	}
	// Build and return a user response with a success code and message
	return buildSuccessfulResponse(newUser.ToResponse(), http.StatusCreated, "user has signed up")
}
