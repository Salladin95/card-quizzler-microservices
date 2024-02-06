package user

import (
	"context"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/lib"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const CollectionPath = "auth"

// getCollection returns a reference to a MongoDB collection in the specified database.
// The collection is derived from the repository's MongoDB configuration and CollectionPath constant.
// It is used to interact with the specified collection in the MongoDB database.
func (repo *repository) getCollection() *mongo.Collection {
	return repo.db.Database(repo.dbCfg.MongoDbName).Collection(CollectionPath)
}

// GetUsers retrieves all users from the MongoDB collection.
// It returns a slice of user.User pointers and an error.
// The users are sorted by the 'createdAt' field in descending order.
// The function utilizes the specified context for the MongoDB operations.
func (repo *repository) GetUsers(ctx context.Context) ([]*user.User, error) {
	collection := repo.getCollection()
	// Set options for sorting by 'createdAt' field in descending order
	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("fetching users", err)
	}
	defer cursor.Close(ctx)
	// Iterate through the cursor and decode users into a slice
	var users []*user.User
	for cursor.Next(ctx) {
		var u user.User
		err := cursor.Decode(&u)
		if err != nil {
			return nil, goErrorHandler.OperationFailure("decoding user into slice", err)
		}
		users = append(users, &u)
	}
	return users, nil
}

// GetById retrieves a user by Id from the MongoDB collection.
// It takes a context and a string representation of the user Id as input.
// The function returns a pointer to a user.User and an error.
// If the user is not found, it returns a custom error with context and wraps the original error.
func (repo *repository) GetById(ctx context.Context, uid string) (*user.User, error) {
	collection := repo.getCollection()
	// Convert the string representation of the user Id to MongoDB ObjectID
	docID, err := toObjectId(uid)
	if err != nil {
		return nil, err
	}
	// Initialize a user.User variable to store the retrieved user
	var user user.User
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&user)
	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by email from the MongoDB collection
func (repo *repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	collection := repo.getCollection()
	// Initialize a user.User variable to store the retrieved user
	var user user.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}
	return &user, nil
}

// CreateUser inserts a new user into the MongoDB collection.
// It takes a context and a SignUpDto as input.
// The function returns a pointer to the created user.User and an error.
// If a user with the same email already exists, it returns a custom error.
func (repo *repository) CreateUser(ctx context.Context, createUserDto userEntities.SignUpDto) (*user.User, error) {
	collection := repo.getCollection()
	// Check if a user with the same email already exists
	_, err := repo.GetByEmail(ctx, createUserDto.Email)
	if err == nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("user with email - %s already exists", createUserDto.Email))
	}
	hashedPassword, err := lib.HashPassword(createUserDto.Password)
	if err != nil {
		return nil, err
	}
	user := user.User{
		ID:        uuid.New(),
		Name:      createUserDto.Name,
		Email:     createUserDto.Email,
		Password:  hashedPassword,
		Birthday:  createUserDto.Birthday,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("inserting user", err)
	}
	return &user, nil
}

// UpdateUser updates a user in the MongoDB collection based on the provided user Id (uid).
// It takes a context, a string representation of the user Id, and an UpdateDto as input.
// The function returns a pointer to the updated user.User and an error.
// If the user is not found by the provided user Id, it returns an error.
// If specific fields (Email, Name, Password) are provided in the UpdateDto,
// the corresponding fields in the user document will be updated.
func (repo *repository) UpdateUser(ctx context.Context, uid string, updateUserDto userEntities.UpdateDto) (*user.User, error) {
	collection := repo.getCollection()
	// Convert the string representation of the user Id to MongoDB ObjectID
	docID, err := toObjectId(uid)
	if err != nil {
		return nil, err
	}
	user, err := repo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	// Update user fields if corresponding values are provided in the UpdateDto
	if updateUserDto.Email != "" {
		user.Email = updateUserDto.Email
	}
	if updateUserDto.Name != "" {
		user.Name = updateUserDto.Name
	}
	if updateUserDto.Password != "" {
		user.Password = updateUserDto.Password
	}

	// Define the filter to identify the user document to update
	filter := bson.M{"_id": docID}
	// Define the update with $set operator for the modified fields
	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"email":     user.Email,
			"password":  user.Password,
			"updatedAt": time.Now(),
		},
	}
	err = collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return nil, goErrorHandler.OperationFailure("updating user", err)
	}
	return user, nil
}

// DeleteUser removes a user from the MongoDB collection
func (repo *repository) DeleteUser(ctx context.Context, uid string) error {
	collection := repo.getCollection()
	result, err := collection.DeleteOne(ctx, bson.M{"id": uid})
	if err != nil {
		return goErrorHandler.OperationFailure("deleting user", err)
	}
	if result.DeletedCount == 0 {
		return goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}
	return nil
}
