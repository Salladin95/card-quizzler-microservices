package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/lib"
	userService "github.com/Salladin95/card-quizzler-microservices/api-service/user"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// pushToQueue pushes data to a named queue in the AMQP exchange.
func (bh *brokerHandlers) pushToQueue(ctx context.Context, name string, data []byte) error {
	// Create a new EventEmitter for the specified AMQP exchange.
	emitter, err := rmqtools.NewEventEmitter(bh.rabbit, AmqpExchange)
	if err != nil {
		return goErrorHandler.OperationFailure("create event emitter", err)
	}

	// Push the data to the named queue using the EventEmitter.
	err = emitter.Push(ctx, name, data)
	if err != nil {
		return goErrorHandler.OperationFailure("push event", err)
	}

	return nil
}

// pushToQueueFromEndpoint handles HTTP requests to push data to a named queue from an endpoint.
func (bh *brokerHandlers) pushToQueueFromEndpoint(c echo.Context, key string) error {
	var requestDTO interface{}

	// Read the request body and unmarshal it into the corresponding DTO.
	if err := c.Bind(&requestDTO); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Marshal the DTO struct into JSON.
	marshalledDto, err := json.Marshal(requestDTO)
	if err != nil {
		return goErrorHandler.OperationFailure("marshal dto", err)
	}

	// Push the marshalled DTO to the named queue using pushToQueue method.
	err = bh.pushToQueue(c.Request().Context(), key, marshalledDto)
	if err != nil {
		return err
	}

	// Respond with a success message indicating the push operation.
	c.String(http.StatusOK, fmt.Sprintf("PUSHED TO QUEUE FROM %s", key))
	return nil
}

// GenerateTokenPair generates a pair of JWTs: access token and refresh token.
// Parameters:
// - name: User's name.
// - email: User's email.
// - id: User's UUID.
// - cfg: JWT configuration.
// Returns a pointer to the generated token pair or an error.
func GenerateTokenPair(name string, email string, id uuid.UUID, cfg config.JwtCfg) (*entities.TokenPair, error) {
	// Generate an access token with a short expiration time
	at, err := GenerateToken(name, email, id, cfg.AccessTokenExpTime, cfg.JWTAccessSecret)
	if err != nil {
		return nil, err
	}

	// Generate a refresh token with a longer expiration time
	rt, err := GenerateToken(name, email, id, cfg.RefreshTokenExpTime, cfg.JWTRefreshSecret)
	if err != nil {
		return nil, err
	}

	// Create a token pair with the generated tokens
	return &entities.TokenPair{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

// GenerateToken generates a JWT with user information.
// Parameters:
// - name: User's name.
// - email: User's email.
// - id: User's UUID.
// - exp: Token expiration duration.
// - secret: Secret key used for signing the token.
// Returns the generated JWT or an error.
func GenerateToken(name string, email string, id uuid.UUID, exp time.Duration, secret string) (string, error) {
	// Create JWT claims with user information
	claims := entities.GetJwtUserClaims(name, email, id, exp)

	// Create a new JWT token with the specified signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the provided secret key
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// If signing fails, return an error with appropriate context
		return "", goErrorHandler.NewError(goErrorHandler.ErrInternalFailure, fmt.Errorf("failed to generate token"))
	}
	// Return the generated JWT
	return t, nil
}

// SetHttpOnlyCookie sets an HTTP-only cookie in the given Echo context.
func SetHttpOnlyCookie(c echo.Context, name, value string, exp time.Duration, path string) {
	// Create a new HTTP cookie
	cookie := new(http.Cookie)

	// Set cookie properties
	cookie.Name = name                     // Cookie name
	cookie.Value = value                   // Cookie value
	cookie.Expires = time.Now().Add(exp)   // Cookie expiration time
	cookie.HttpOnly = true                 // Set the cookie as HTTP-only
	cookie.Path = path                     // Cookie path
	cookie.SameSite = http.SameSiteLaxMode // Set SameSite policy (Lax mode)
	cookie.Secure = true                   // Set the cookie as secure (HTTPS only)

	// Set the cookie in the Echo context
	c.SetCookie(cookie)
}

func buildResponse(c echo.Context, res *userService.Response, unmarshalTo interface{}) error {
	code := res.GetCode()
	if code >= http.StatusBadRequest {
		return c.JSON(int(res.GetCode()), entities.JsonResponse{Message: res.GetMessage()})
	}
	err := lib.UnmarshalData(res.GetData(), &unmarshalTo)
	if err != nil {
		return err
	}
	return c.JSON(int(res.GetCode()), entities.JsonResponse{Message: res.GetMessage(), Data: unmarshalTo})
}
