package cachedRepository

import (
	"context"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/constants"
	appEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/lib"
	userEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	userRepo "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/repository"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"github.com/go-redis/redis"
	"time"
)

// cachedRepository is a repository implementation that caches data using Redis.
type cachedRepository struct {
	broker      rmqtools.MessageBroker // Message broker interface for communication
	redisClient *redis.Client          // Redis client used for caching
	repo        userRepo.Repository    // Underlying repository for fetching data
	exp         time.Duration          // Expiration time for cached data
}

// CachedRepository implements the CachedRepository interface for user data operations with read and write through caching pattern.
// It attempts to read data from cache first. If the data is found in the cache, it returns it.
// If the cache read fails, it fetches the data from the repository, updates the cache, and returns the data.
// On data mutation, it updates the mutated data in the cache.
type CachedRepository interface {
	userRepo.Repository
	GetEmailVerificationCode(ctx context.Context, email string) (*appEntities.EmailCode, error)
	SetEmailVerificationCode(ctx context.Context, payload []byte) error
}

// NewCachedUserRepo creates a new instance of CachedRepository.
func NewCachedUserRepo(
	broker rmqtools.MessageBroker, // Message broker interface for communication
	redisClient *redis.Client, // Redis client used for caching
	userRep userRepo.Repository, // Underlying repository for fetching data
) CachedRepository {
	return &cachedRepository{
		broker:      broker,           // Assign message broker
		redisClient: redisClient,      // Assign Redis client
		repo:        userRep,          // Assign underlying repository
		exp:         60 * time.Minute, // Set expiration time for cached data
	}
}

// GetUsers retrieves user data either from the cache or the underlying repository.
// It first attempts to read users from the cache. If successful, it returns the cached users.
// If reading from the cache fails (cache miss), it fetches users from the underlying repository,
// caches the fetched users, and publishes an event to RabbitMQ indicating that users were fetched.
// It returns the fetched users or an error if fetching users from the repository fails.
func (cr *cachedRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	var users []*user.User
	// Try to read users from the cache
	err := cr.readCacheByHashedKey(&users, userRootKey, usersKey)
	if err != nil {
		// If cache read fails, fetch users from the underlying repository
		users, err = cr.repo.GetUsers(ctx)
		if err != nil {
			return nil, err // Return error if fetching users from the repository fails
		}

		// Cache the fetched users
		cr.setCacheInPipeline(userRootKey, usersKey, users)
		// Log message indicating users were retrieved from the repository and cached
		cr.log(
			ctx,
			"users retrieved from repository and cached",
			"info",
			"GetUsers",
		)
	}

	// Publish an event to RabbitMQ indicating that users were fetched
	cr.broker.PushToQueue(ctx, constants.FetchedUsersKey, users)
	return users, nil // Return the fetched users
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
	err := cr.readCacheByHashedKey(&user, userRootKey, cr.userHashKey(uid))
	if err != nil {
		// If cache read fails, fetch the user from the underlying repository
		user, err = cr.repo.GetById(ctx, uid)
		if err != nil {
			return nil, err
		}

		// Cache the fetched user using both the hash key and email as cache keys
		cr.setCacheInPipeline(userRootKey, cr.userHashKey(uid), user)
		cr.setCacheInPipeline(userRootKey, cr.userHashKey(user.Email), user)

		// Log message indicating user was retrieved from the repository and cached
		cr.log(
			ctx,
			"user retrieved from repository and cached",
			"info",
			"GetById",
		)
	}

	// Publish an event to RabbitMQ indicating that the user was fetched
	cr.broker.PushToQueue(ctx, constants.FetchedUserKey, user)

	return user, nil
}

// GetByEmail retrieves a user by their email, either from cache or the underlying repository.
func (cr *cachedRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user *user.User
	// Try to read user from cache using the email as cache key
	err := cr.readCacheByKey(&user, email)
	if err != nil {
		// If cache read fails, fetch user from the underlying repository
		user, err = cr.repo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		// Cache the fetched user using both the hash key derived from user ID and email as cache keys
		cr.setCacheInPipeline(userRootKey, cr.userHashKey(user.ID.String()), user)
		cr.setCacheInPipeline(userRootKey, cr.userHashKey(user.Email), user)

		// Log message indicating user was retrieved from the repository and cached
		cr.log(
			ctx,
			"user retrieved from repository and cached",
			"info",
			"GetByEmail",
		)
	}

	cr.broker.PushToQueue(ctx, constants.FetchedUserKey, user)
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
		return nil, err // Return error if creating the user fails
	}

	// Cache the newly created user using both the hash key derived from the user ID and the email as cache keys
	cr.setCacheInPipeline(userRootKey, cr.userHashKey(createdUser.ID.String()), createdUser)
	cr.setCacheInPipeline(userRootKey, cr.userHashKey(createdUser.Email), createdUser)

	// Clear the cache for the user list
	cr.clearCacheByKeys(userRootKey, usersKey)

	// Publish an event to RabbitMQ indicating that the user was created
	cr.broker.PushToQueue(ctx, constants.CreatedUserKey, createdUser)

	// Log message indicating a new user creation. User data has been cached using both ID and EMAIL keys. Additionally, the cache for the USERS key has been reset, and an event has been generated.
	cr.log(
		ctx,
		"New user created. User data cached by ID & EMAIL. Cache reset for USERS key. Event generated.",
		"info",
		"CreateUser",
	)

	return createdUser, nil // Return the newly created user
}

// UpdateUser updates an existing user with the provided data.
// It first updates the user using the underlying repository.
// If successful, it caches the updated user using both the hash key derived from the user ID
// and the email as cache keys. Additionally, it publishes an event to RabbitMQ indicating that the user was updated.
// It then clears the cache for the user list.
// It returns the updated user or an error if updating the user fails.
func (cr *cachedRepository) UpdateUser(
	ctx context.Context,
	uid string,
	updateUserDto userEntities.UpdateUserDto,
) (*user.User, error) {
	// Update the user using the underlying repository
	updatedUser, err := cr.repo.UpdateUser(ctx, uid, updateUserDto)
	if err != nil {
		return nil, err
	}

	// Clear the cache
	cr.clearCacheByKey(userRootKey)

	// Publish an event to RabbitMQ indicating that the user was updated
	cr.broker.PushToQueue(ctx, constants.UpdatedUserKey, updatedUser)

	// Log message indicating a new user creation. User data has been cached using both ID and EMAIL keys. Additionally, the cache for the USERS key has been reset, and an event has been generated.
	cr.log(
		ctx,
		"User updated. User cache updated for keys ID & EMAIL. Cache reset for USERS key. Event generated.",
		"info",
		"UpdateEmail",
	)

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

	// Clear the cache
	cr.clearCacheByKey(userRootKey)

	// Publish an event to RabbitMQ indicating that the user was deleted
	cr.broker.PushToQueue(ctx, constants.DeletedUserKey, u)

	// Log message indicating a new user creation. User data has been cached using both ID and EMAIL keys. Additionally, the cache for the USERS key has been reset, and an event has been generated.
	cr.log(
		ctx,
		"User deleted. User cache reset for keys ID, EMAIL, USERS. Event generated.",
		"info",
		"DeleteUser",
	)
	return nil
}

// SetEmailVerificationCode sets the email verification code for the specified email in the cache.
// It first unmarshals the payload into an EmailCode struct and verifies its validity.
// If successful, it sets the email verification code in the cache with a 2-minute expiration.
// It logs the action and returns nil if successful, otherwise it returns an error.
func (cr *cachedRepository) SetEmailVerificationCode(ctx context.Context, payload []byte) error {
	// Unmarshal payload into EmailCode struct
	var emailCode appEntities.EmailCode
	err := lib.UnmarshalData(payload, &emailCode)
	if err != nil {
		return err
	}

	// Verify email code
	err = emailCode.Verify()
	if err != nil {
		cr.log(ctx, err.Error(), "error", "SetEmailVerificationCode")
		return err
	}

	// Set email verification code in cache with 2-minute expiration
	err = cr.redisClient.Set(cr.codeKey(emailCode.Email), payload, 2*time.Minute).Err()

	if err != nil {
		cr.log(ctx, err.Error(), "error", "SetEmailVerificationCode")
		return goErrorHandler.OperationFailure("set cache", err)
	}

	// Log action
	cr.log(
		ctx,
		"Email verification code is saved to cache",
		"info",
		"SetEmailVerificationCode",
	)

	return nil
}

// GetEmailVerificationCode retrieves the email verification code for the specified email from the cache.
// It reads the email code from the cache using the email as the key.
// If successful, it verifies the email code and returns it along with nil error.
// It logs the action and returns an error if the code is invalid or if retrieval from the cache fails.
func (cr *cachedRepository) GetEmailVerificationCode(ctx context.Context, email string) (*appEntities.EmailCode, error) {
	// Initialize EmailCode struct
	var emailCode appEntities.EmailCode

	// Read email code from cache
	err := cr.readCacheByKey(&emailCode, cr.codeKey(email))
	if err != nil {
		cr.log(ctx, err.Error(), "error", "GetEmailVerificationCode")
		return nil, err
	}

	// Verify email code
	err = emailCode.Verify()
	if err != nil {
		cr.log(ctx, err.Error(), "error", "GetEmailVerificationCode")
		return nil, err
	}

	// Log action
	cr.log(
		ctx,
		"Email verification code has been extracted from cache",
		"info",
		"SetEmailVerificationCode",
	)

	return &emailCode, nil
}

// log sends a log message to the message broker.
func (cr *cachedRepository) log(ctx context.Context, message, level, method string) {
	var log appEntities.LogMessage // Create a new LogMessage struct
	// Push log message to the message broker
	cr.broker.PushToQueue(
		ctx,
		constants.LogCommand, // Specify the log command constant
		// Generate log message with provided details
		log.GenerateLog(message, level, method, "cached repository"),
	)
}
