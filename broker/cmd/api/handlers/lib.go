package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	auth "github.com/Salladin95/card-quizzler-microservices/broker-service/auth"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/config"
	"github.com/Salladin95/card-quizzler-microservices/broker-service/cmd/api/entities"
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

// Unmarshal user and returns auth response
func buildUserResponse(c echo.Context, res *auth.Response) error {
	resCode := int(res.GetCode())
	if resCode >= 400 {
		return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage()})
	}
	var userResponse entities.UserResponse
	err := json.Unmarshal(res.GetData(), &userResponse)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}
	return c.JSON(resCode, entities.JsonResponse{Message: res.GetMessage(), Data: userResponse})
}

type tokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func GenerateTokenPair(name string, email string, id uuid.UUID, cfg config.JwtCfg) (*tokenPair, error) {
	//at, err := GenerateToken(name, email, id, cfg.AccessTokenExpTime, cfg.JWTAccessSecret)
	at, err := GenerateToken(name, email, id, 15*time.Second, cfg.JWTAccessSecret)
	if err != nil {
		return nil, err
	}
	//rt, err := GenerateToken(name, email, id, cfg.RefreshTokenExpTime, cfg.JWTRefreshSecret)
	rt, err := GenerateToken(name, email, id, 1*time.Minute, cfg.JWTRefreshSecret)
	if err != nil {
		return nil, err
	}
	return &tokenPair{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

func GenerateToken(name string, email string, id uuid.UUID, exp time.Duration, secret string) (string, error) {
	claims := entities.GetJwtUserClaims(name, email, id, exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", goErrorHandler.NewError(goErrorHandler.ErrInternalFailure, fmt.Errorf("failed to generate token"))
	}
	return t, nil
}

func SetHttpOnlyCookie(c echo.Context, name, value string, exp time.Duration, path string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(exp)
	cookie.HttpOnly = true
	cookie.Path = path
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Secure = true
	c.SetCookie(cookie)
}
