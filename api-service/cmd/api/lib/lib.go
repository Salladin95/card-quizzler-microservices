package lib

import (
	"encoding/json"
	"errors"
	"github.com/Salladin95/card-quizzler-microservices/api-service/cmd/api/config"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// DataWithVerify is an interface that requires a Verify method.
type DataWithVerify interface {
	Verify() error
}

// Verify validates the given structure
func Verify(data interface{}) error {
	// Create a new validator instance.
	validate := validator.New()

	// Validate the SignUpDto structure.
	if err := validate.Struct(data); err != nil {
		// Convert validation errors and return a ValidationFailure error.
		return goErrorHandler.ValidationFailure(goErrorHandler.ConvertValidationErrors(err))
	}

	return nil
}

// BindBody bind body from request to bindTo param
// Note - bindTo must be pointer !!!
func BindBody(c echo.Context, bindTo interface{}) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(&bindTo); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}
	return nil
}

// BindBodyAndVerify binds the request body to a DataWithVerify interface
// and then calls the Verify method on the provided data.
// Note - bindTo must be pointer !!!
func BindBodyAndVerify(c echo.Context, bindTo DataWithVerify) error {
	// Bind the request body to the DataWithVerify interface
	if err := c.Bind(bindTo); err != nil {
		return goErrorHandler.BindRequestToBodyFailure(err)
	}

	// Call the Verify method on the provided bindTo
	err := bindTo.Verify()
	return err
}

// UnmarshalData unmarshals JSON data into the provided unmarshalTo interface.
// It returns an error if any issues occur during the unmarshaling process.
// Note - unmarshalTo must be pointer !!!
func UnmarshalData(data []byte, unmarshalTo interface{}) error {
	err := json.Unmarshal(data, unmarshalTo)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal data", err)
	}
	return nil
}

// MarshalData marshals data into a JSON-encoded byte slice.
// It returns the marshalled data []byte and an error if any issues occur during the marshaling process.
func MarshalData(data interface{}) ([]byte, error) {
	marshalledData, err := json.Marshal(data)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("marshal data", err)
	}
	return marshalledData, nil
}

type services struct {
	Redis  *redis.Client
	Rabbit *amqp091.Connection
}

func InitializeServices(cfg config.AppCfg) services {
	// Connect to RabbitMQ server using the provided URL.
	rabbitConn, err := rmqtools.ConnectToRabbit(cfg.RabbitUrl)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Establish a Redis connection
	redisConn := connectToRedis(cfg.RedisUrl)
	return services{
		Redis:  redisConn,
		Rabbit: rabbitConn,
	}
}

// connectToRedis establishes a connection to a Redis server and returns a Redis client.
// It takes the address of the Redis server as a parameter.
func connectToRedis(addr string) *redis.Client {
	// Create a new Redis client with specified options
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		WriteTimeout: 5 * time.Second, // Maximum time to wait for write operations
		ReadTimeout:  5 * time.Second, // Maximum time to wait for read operations
		DialTimeout:  3 * time.Second, // Maximum time to wait for a connection to be established
		MaxRetries:   3,               // Maximum number of retries before giving up on a command
	})
}

type JwtUser struct {
	Id string `json:"id"`
}

type JwtUserClaims struct {
	JwtUser
	jwt.RegisteredClaims
}

func GetJwtUserClaims(id string, exp time.Duration) JwtUserClaims {
	return JwtUserClaims{
		JwtUser{id},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}

// GetJwtConfig returns a configuration for JWT authentication.
// It takes a key as a parameter for the signing key.
func GetJwtConfig(key string) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(key), // Set the signing key for JWT
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtUserClaims) // Create new JWT claims using the JwtUserClaims struct
		},
		ErrorHandler: func(c echo.Context, err error) error {
			// HandleEvent errors during JWT authentication and return an unauthorized error
			return goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
		},
	}
}

// ValidateTokenString validates a JWT token string using the provided secret key.
// It returns the decoded claims if the token is valid, otherwise, returns an error.
func ValidateTokenString(tokenString string, secret string) (*JwtUserClaims, error) {
	// Get JWT configuration with the provided secret key
	config := GetJwtConfig(secret)

	// Parse and validate the token with custom claims and signing key
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtUserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return config.SigningKey, nil
		})

	if err != nil || !token.Valid {
		// Log and return an unauthorized error if the token is invalid
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	// Extract claims from the token
	claims, ok := token.Claims.(*JwtUserClaims)
	if !ok {
		// Log and return an unauthorized error if claims are invalid
		return nil, goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	return claims, nil
}

// ClearCookies removes refresh-token cookie from the provided echo.Context.
// It takes the echo.Context as parameter.
func ClearCookies(c echo.Context) {
	// Reset the specified cookies
	ResetCookies(c, "refresh-token")
}

// ResetCookies clears cookies by setting new cookies with an expiration time in the past.
// It takes echo.Context and cookie names as parameters.
func ResetCookies(c echo.Context, cookieNames ...string) {
	for _, cookieName := range cookieNames {
		// Create a new cookie with an expiration time in the past
		cookie := new(http.Cookie)
		cookie.Name = cookieName
		cookie.Value = ""
		cookie.Expires = time.Now().Add(-24 * time.Hour) // Set expiration time to the past
		// Set the cookie in the response
		c.SetCookie(cookie)
	}
}

// ExtractAccessToken extracts the access token from the Authorization header in the echo.Context.
// It returns the extracted access token and an error, if any.
func ExtractAccessToken(c echo.Context) (string, error) {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")

	// Split the Authorization header to retrieve the token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// If the token is missing or has an incorrect format, clear cookies and return an error
		ClearCookies(c)
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, errors.New("token is missing or has an incorrect format"))
	}

	// Extract and return the Bearer token
	bearerToken := strings.TrimPrefix(tokenParts[1], "Bearer ")
	bearerToken = strings.Trim(bearerToken, "\"")
	return bearerToken, nil
}

// ExtractRefreshToken extracts the refresh token from the "refresh-token" cookie in the echo.Context.
// It returns the extracted refresh token and an error, if any.
func ExtractRefreshToken(c echo.Context) (string, error) {
	// Get the "refresh-token" cookie from the request
	tokenCookie, err := c.Cookie("refresh-token")
	if err != nil || tokenCookie.Value == "" {
		// If the refresh token cookie is missing or empty, clear cookies and return an error
		ClearCookies(c)
		return "", goErrorHandler.NewError(goErrorHandler.ErrUnauthorized, err)
	}

	// Extract and return the refresh token
	tokenString := tokenCookie.Value
	return tokenString, nil
}

func LogInfo(msg string, args ...any) {
	slog.Info(msg, args...)
}

func LogError(err error, args ...any) {
	slog.Error(err.Error(), args...)
}

func LogRequestInfo(c echo.Context, msg string) {
	slog.Info(
		msg,
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)
}

func LogRequestError(c echo.Context, msg string) {
	slog.Error(
		msg,
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)
}
