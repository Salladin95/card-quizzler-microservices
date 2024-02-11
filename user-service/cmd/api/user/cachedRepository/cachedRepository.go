package cachedRepository

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	userEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	userRepo "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/repository"
	"github.com/go-redis/redis"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type cachedRepository struct {
	redisClient *redis.Client       // redisClient is the Redis client used for caching.
	rabbitConn  *amqp091.Connection // rabbitConn is the AMQP connection used for firing events.
	repo        userRepo.Repository // repo is the underlying repository from which to fetch data.
	exp         time.Duration       // exp is the expiration time for cached data.
}

func NewCachedUserRepo(
	redisClient *redis.Client,
	rabbitConn *amqp091.Connection,
	userRep userRepo.Repository,
) CachedRepository {
	return &cachedRepository{
		redisClient: redisClient,
		repo:        userRep,
		exp:         60 * time.Minute,
		rabbitConn:  rabbitConn,
	}
}

type LogMessage struct {
	FromService string `json:"fromService" validate:"required"`
	Message     string `json:"message" validate:"required"`
	Level       string `json:"level" validate:"required"`
	Name        string `json:"name" validate:"omitempty"`
	Method      string `json:"method" validate:"omitempty"`
}

// GetUsers retrieves user data either from the cache or the underlying repository.
// It first attempts to read users from the cache. If successful, it returns the cached users.
// If reading from the cache fails (cache miss), it fetches users from the underlying repository,
// caches the fetched users, and publishes an event to RabbitMQ indicating that users were fetched.
// It returns the fetched users or an error if fetching users from the repository fails.
func (cr *cachedRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	var users []*user.User
	// Try to read users from the cache
	err := cr.readCacheByKey(&users, usersKey)
	if err != nil {
		// If cache read fails, fetch users from the underlying repository
		users, err = cr.repo.GetUsers(ctx)
		if err != nil {
			return nil, err
		}

		// Cache the fetched users
		cr.setCacheByKey(usersKey, users)
		cr.pushToQueue(
			ctx,
			constants.LogCommand,
			generateLog("users retrieved from repo and set as a cache", "info", "GetUsers"),
		)
	}

	// Publish an event to RabbitMQ indicating that users were fetched
	cr.pushToQueue(ctx, constants.FetchedUsersKey, users)
	return users, nil
}

// GetById retrieves a user by their ID, either from the cache or the underlying repository.
// It first attempts to read the user from the cache using the hash key derived from the user ID.
// If successful, it returns the cached user.
// If reading from the cache fails (cache miss), it fetches the user from the underlying repository
// and caches the fetched user using both the hash key and email as cache keys.
// Additionally, it publishes an event to RabbitMQ indicating that the user was fetched.
// It returns the fetched user or an error if fetching the user from the repository fails.
func (cr *cachedRepository) GetById(ctx context.Context, uid string) (*user.User, error) {
	var user *user.User

	// Try to read the user from the cache using the hash key derived from the user ID
	err := cr.readCacheByKey(&user, cr.userHashKey(uid))
	if err != nil {
		// If cache read fails, fetch the user from the underlying repository
		user, err = cr.repo.GetById(ctx, uid)
		if err != nil {
			return nil, err
		}

		// Cache the fetched user using both the hash key and email as cache keys
		cr.setCacheByKey(cr.userHashKey(uid), user)
		cr.setCacheByKey(user.Email, user)
	}

	// Publish an event to RabbitMQ indicating that the user was fetched
	cr.pushToQueue(ctx, constants.FetchedUserKey, user)

	// If cache read succeeds, return the user from the cache
	log.Println("[user-service] User has been retrieved")
	return user, nil
}

// GetByEmail retrieves a user by their email, either from cache or the underlying repository.
func (cr *cachedRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user *user.User
	// Try to read user from cache using the email as cache key
	err := cr.readCacheByKey(&user, email)
	if err != nil {
		// If cache read fails, fetch user from the underlying repository
		user, err := cr.repo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		// Cache the fetched user using both the hash key derived from user ID and email as cache keys
		cr.setCacheByKey(cr.userHashKey(user.ID.String()), user)
		cr.setCacheByKey(email, user)
	}

	cr.pushToQueue(ctx, constants.FetchedUserKey, user)
	// If cache read succeeds, return user from cache
	log.Println("[user-service] User has been extracted from cache")
	return user, nil
}

// CreateUser creates a new user using the provided user data.
// It first creates the user using the underlying repository.
// If successful, it caches the newly created user using both the hash key derived from the user ID
// and the email as cache keys. Additionally, it publishes an event to RabbitMQ indicating that the user was created.
// It then clears the cache for the user list.
// It returns the newly created user or an error if creating the user fails.
func (cr *cachedRepository) CreateUser(
	ctx context.Context,
	createUserDto userEntities.SignUpDto,
) (*user.User, error) {
	// Create the user using the underlying repository
	createdUser, err := cr.repo.CreateUser(ctx, createUserDto)
	if err != nil {
		return nil, err
	}

	// Cache the newly created user using both the hash key derived from the user ID and the email as cache keys
	cr.setCacheByKey(cr.userHashKey(createdUser.ID.String()), createdUser)
	cr.setCacheByKey(createdUser.Email, createdUser)

	// Publish an event to RabbitMQ indicating that the user was created
	cr.pushToQueue(ctx, constants.CreatedUserKey, createdUser)

	// Clear the cache for the user list
	cr.clearCacheByKey(usersKey)

	return createdUser, nil
}

// UpdateUser updates an existing user with the provided data.
// It first updates the user using the underlying repository.
// If successful, it caches the updated user using both the hash key derived from the user ID
// and the email as cache keys. Additionally, it publishes an event to RabbitMQ indicating that the user was updated.
// It then clears the cache for the user list.
// It returns the updated user or an error if updating the user fails.
func (cr *cachedRepository) UpdateUser(
	ctx context.Context, uid string,
	updateUserDto userEntities.UpdateDto,
) (*user.User, error) {
	// Update the user using the underlying repository
	updatedUser, err := cr.repo.UpdateUser(ctx, uid, updateUserDto)
	if err != nil {
		return nil, err
	}

	// Cache the updated user using both the hash key derived from user ID and the email as cache keys
	cr.setCacheByKey(cr.userHashKey(updatedUser.ID.String()), updatedUser)
	cr.setCacheByKey(updatedUser.Email, updatedUser)

	// Publish an event to RabbitMQ indicating that the user was updated
	cr.pushToQueue(ctx, constants.UpdatedUserKey, updatedUser)

	// Clear the cache for the user list
	cr.clearCacheByKey(usersKey)

	return updatedUser, nil
}

// DeleteUser deletes a user with the specified ID.
// It first retrieves the user to clear its cache.
// If successful, it deletes the user using the underlying repository.
// Additionally, it clears the cache for the deleted user and the user list.
// It also publishes an event to RabbitMQ indicating that the user was deleted.
// It returns nil if the operation succeeds, otherwise it returns an error.
func (cr *cachedRepository) DeleteUser(ctx context.Context, uid string) error {
	// Retrieve the user to clear its cache
	u, err := cr.GetById(ctx, uid)

	// Delete the user using the underlying repository
	err = cr.repo.DeleteUser(ctx, uid)
	if err != nil {
		return err
	}

	// Clear the cache for the deleted user and the user list
	cr.clearCacheByKey(cr.userHashKey(uid))
	cr.clearCacheByKey(u.Email)
	cr.clearCacheByKey(usersKey)

	// Publish an event to RabbitMQ indicating that the user was deleted
	cr.pushToQueue(ctx, constants.DeletedUserKey, u)
	return nil
}

func generateLog(message string, level string, method string) LogMessage {
	return LogMessage{
		Level:       level,
		Method:      method,
		FromService: "user-service",
		Message:     message,
		Name:        "working with cachedRepository",
	}
}
