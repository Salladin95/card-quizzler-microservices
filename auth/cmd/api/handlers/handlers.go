package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/constants"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	repo "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/repository"
	auth "github.com/Salladin95/card-quizzler-microservices/auth-service/proto"
	"log"
	"net/http"
)

// AuthServer is the gRPC server implementation for authentication-related operations.
type AuthServer struct {
	auth.UnimplementedAuthServer                 // Embed the autogenerated UnimplementedAuthServer to satisfy the interface.
	Repo                         repo.Repository // User repository
}

// SignIn is a gRPC service method that processes a sign-in request.
// It takes a context and a gRPC request (auth.SignInRequest) as input.
// The function returns a gRPC response (auth.Response) and an error.
// The response includes a status code and a message.
func (authServer *AuthServer) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.Response, error) {
	// Print a message indicating the start of processing the sign-in request
	fmt.Println("******* auth service - start processing signin request ********")
	// Extract payload from the gRPC request
	reqPayload := req.GetPayload()
	// Create a SignInDto from the request payload
	signInDto := user.SignInDto{Email: reqPayload.Email, Password: reqPayload.Password}
	// Verify the SignInDto structure
	err := signInDto.Verify()
	if err != nil {
		// Return a response with the mapped error status code and message if verification fails
		return &auth.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Fetch user by email from the repository
	fetchedUser, err := authServer.Repo.GetByEmail(ctx, signInDto.Email)
	if err != nil {
		// Return a response with the mapped error status code and message if fetching user fails
		return &auth.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Check if the fetched user's email matches the sign-in email and verify the password
	isPasswordInvalid := authServer.Repo.CompareHashAndPassword(fetchedUser.Password, signInDto.Password)
	if fetchedUser.Email != signInDto.Email || isPasswordInvalid != nil {
		// Return a response with the mapped error status code and message if authentication fails
		return &auth.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Build and return a user response with a success code and message
	return buildUserResponse(fetchedUser, http.StatusOK, "user has signed in")
}

// SignUp is a gRPC service method that handles the user sign-up request.
// It takes a context and a gRPC request (auth.SignUpRequest) as input.
// The function returns a gRPC response (auth.Response) and an error.
// The response includes a status code and a message.
func (authServer *AuthServer) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.Response, error) {
	// Print a message indicating the start of processing the sign-up request
	fmt.Println("******* auth service - start processing signUp request ********")
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
		return &auth.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Create a new user by calling the CreateUser method in the repository
	newUser, err := authServer.Repo.CreateUser(ctx, signUpDto)
	if err != nil {
		// Return a response with the mapped error status code and message if user creation fails
		return &auth.Response{Code: getErrorStatus(err), Message: getErrorMessage(err)}, nil
	}
	// Build and return a user response with a success code and message
	return buildUserResponse(newUser, http.StatusCreated, "user has signed up")
}

func HandleRabbitPayload(key string, payload []byte) {
	fmt.Print("START PROCESSING MESSAGE")
	switch key {
	case constants.SignInKey:
		fmt.Printf("******* - %v\n\n", payload)
		var signInDto user.SignInDto
		err := json.Unmarshal(payload, &signInDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN IN *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %s\n", key)
		fmt.Printf("payload - %v\n\n", signInDto)
	case constants.SignUpKey:
		var signUpDto user.SignUpDto
		err := json.Unmarshal(payload, &signUpDto)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("******************* SIGN UP *****************")
		fmt.Printf("MESSAGE FROM QUEUE - %v", key)
	default:
		log.Panic("handlePayload: unknown payload name")
	}
}
