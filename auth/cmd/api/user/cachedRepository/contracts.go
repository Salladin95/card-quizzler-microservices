package cachedRepository

import userRepo "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/repository"

// CachedRepository implements the CachedRepository interface for user data operations with read and write through caching pattern.
// It attempts to read data from cache first. If the data is found in the cache, it returns it.
// If the cache read fails, it fetches the data from the repository, updates the cache, and returns the data.
// On data mutation, it updates the mutated data in the cache.
type CachedRepository interface {
	userRepo.Repository
}
