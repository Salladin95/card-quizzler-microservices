package cachedRepository

import (
	"context"
	userEntities "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/model"
	userRepo "github.com/Salladin95/card-quizzler-microservices/user-service/cmd/api/user/repository"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type cachedRepository struct {
	redisClient *redis.Client       // redisClient is a Redis client used for caching
	repo        userRepo.Repository // repo is the underlying repository from which to fetch data
	userKey     string              // userKey is the key used to store/retrieve user data from cache
	exp         time.Duration       // exp is the expiration time for cached data
}

// GetUsers retrieves user data either from cache or the underlying repository.
func (cr *cachedRepository) GetUsers(ctx context.Context) ([]*user.User, error) {
	var cachedUsers []*user.User
	// Try to read users from cache
	err := cr.readCacheByKey(&cachedUsers, cr.userKey)
	if err != nil {
		// If cache read fails, fetch users from the underlying repository
		users, err := cr.repo.GetUsers(ctx)
		if err != nil {
			return nil, err
		}
		// Cache the fetched users
		cr.SetCacheByKey(cr.userKey, users)
		log.Println("users retrieved from repository")
		return users, nil
	}
	// If cache read succeeds, return users from cache
	log.Println("users retrieved from cache")
	return cachedUsers, nil
}

// GetById retrieves a user by their ID, either from cache or the underlying repository.
func (cr *cachedRepository) GetById(ctx context.Context, uid string) (*user.User, error) {
	var cachedUser *user.User
	// Try to read user from cache using the hash key derived from the user ID
	err := cr.readCacheByKey(&cachedUser, cr.userHashKey(uid))
	if err != nil {
		// If cache read fails, fetch user from the underlying repository
		user, err := cr.repo.GetById(ctx, uid)
		if err != nil {
			return nil, err
		}
		// Cache the fetched user using both the hash key and email as cache keys
		cr.SetCacheByKey(cr.userHashKey(uid), user)
		cr.SetCacheByKey(user.Email, user)
		log.Println("user retrieved from repository")
		return user, nil
	}
	// If cache read succeeds, return user from cache
	log.Printf("user - %v has been extracted from cache\n", cachedUser)
	return cachedUser, nil
}

// GetByEmail retrieves a user by their email, either from cache or the underlying repository.
func (cr *cachedRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var cachedUser *user.User
	// Try to read user from cache using the email as cache key
	err := cr.readCacheByKey(&cachedUser, email)
	if err != nil {
		// If cache read fails, fetch user from the underlying repository
		user, err := cr.repo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		// Cache the fetched user using both the hash key derived from user ID and email as cache keys
		cr.SetCacheByKey(cr.userHashKey(user.ID.String()), user)
		cr.SetCacheByKey(email, user)
		return user, nil
	}
	// If cache read succeeds, return user from cache
	log.Printf("user - %v has been extracted from cache\n", cachedUser)
	return cachedUser, nil
}

// CreateUser creates a new user using the provided user data.
// It caches the newly created user and clears the user list cache.
func (cr *cachedRepository) CreateUser(
	ctx context.Context,
	createUserDto userEntities.SignUpDto,
) (*user.User, error) {
	// Create the user using the underlying repository
	createdUser, err := cr.repo.CreateUser(ctx, createUserDto)
	if err != nil {
		return nil, err
	}
	// Cache the newly created user using both the hash key derived from user ID and email as cache keys
	cr.SetCacheByKey(cr.userHashKey(createdUser.ID.String()), createdUser)
	cr.SetCacheByKey(createdUser.Email, createdUser)
	// Clear the cache for the user list
	cr.clearCacheByKey(cr.userKey)
	return createdUser, nil
}

// UpdateUser updates an existing user with the provided data.
// It caches the updated user and clears the user list cache.
func (cr *cachedRepository) UpdateUser(
	ctx context.Context, uid string,
	updateUserDto userEntities.UpdateDto,
) (*user.User, error) {
	// Update the user using the underlying repository
	updatedUser, err := cr.repo.UpdateUser(ctx, uid, updateUserDto)
	if err != nil {
		return nil, err
	}
	// Cache the updated user using both the hash key derived from user ID and email as cache keys
	cr.SetCacheByKey(cr.userHashKey(updatedUser.ID.String()), updatedUser)
	cr.SetCacheByKey(updatedUser.Email, updatedUser)
	// Clear the cache for the user list
	cr.clearCacheByKey(cr.userKey)
	return updatedUser, nil
}

// DeleteUser deletes a user with the specified ID.
// It retrieves the user first to clear its cache, then deletes the user and clears the user list cache.
func (cr *cachedRepository) DeleteUser(ctx context.Context, uid string) error {
	// Retrieve the user to clear its cache
	u, getUserErr := cr.GetById(ctx, uid)
	// Delete the user using the underlying repository
	err := cr.repo.DeleteUser(ctx, uid)
	if err != nil {
		return err
	}
	// Clear the cache for the deleted user and the user list
	cr.clearCacheByKey(cr.userHashKey(uid))
	cr.clearCacheByKey(cr.userKey)
	// If there was an error while retrieving the user, clear the cache for the user's email
	if getUserErr != nil {
		cr.clearCacheByKey(u.Email)
	}
	return nil
}
