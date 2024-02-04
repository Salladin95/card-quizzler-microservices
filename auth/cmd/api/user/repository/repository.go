package user

import (
	"context"
	"fmt"
	userEntities "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/entities"
	user "github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/user/model"
	"github.com/Salladin95/goErrorHandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const CollectionPath = "auth"

func (repo *repository) getCollection() *mongo.Collection {
	return repo.db.Database(repo.dbCfg.MongoDbName).Collection(CollectionPath)
}

// GetUsers retrieves all users from the MongoDB collection
func (repo *repository) GetUsers(ctx context.Context) ([]*user.User, error) {
	collection := repo.getCollection()

	opts := options.Find()
	opts.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, goErrorHandler.OperationFailure("fetching users", err)
	}
	defer cursor.Close(ctx)

	var users []*user.User
	for cursor.Next(ctx) {
		var user user.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, goErrorHandler.OperationFailure("decoding user into slice", err)
		} else {
			users = append(users, &user)
		}
	}

	return users, nil
}

// GetById retrieves a user by ID from the MongoDB collection
func (repo *repository) GetById(ctx context.Context, uid string) (*user.User, error) {
	collection := repo.getCollection()

	var user user.User

	docID, err := toObjectId(uid)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&user)

	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email from the MongoDB collection
func (repo *repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	collection := repo.getCollection()

	var user user.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrNotFound, err)
	}

	return &user, nil
}

// CreateUser inserts a new user into the MongoDB collection
func (repo *repository) CreateUser(ctx context.Context, createUserDto userEntities.SignUpDto) (*user.User, error) {
	collection := repo.getCollection()

	_, err := repo.GetByEmail(ctx, createUserDto.Email)
	if err == nil {
		return nil, goErrorHandler.NewError(goErrorHandler.ErrBadRequest, fmt.Errorf("user with email - %s already exixsts", createUserDto.Email))
	}

	hashedPassword, err := repo.HashPassword(createUserDto.Password)
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

// UpdateUser updates a user in the MongoDB collection
func (repo *repository) UpdateUser(ctx context.Context, uid string, updateUserDto userEntities.UpdateDto) (*user.User, error) {
	collection := repo.getCollection()

	docID, err := toObjectId(uid)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}

	if updateUserDto.Email != "" {
		user.Email = updateUserDto.Email
	}

	if updateUserDto.Name != "" {
		user.Name = updateUserDto.Name
	}

	if updateUserDto.Password != "" {
		user.Password = updateUserDto.Password
	}

	filter := bson.M{"_id": docID}
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

func toObjectId(id string) (primitive.ObjectID, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return docID, goErrorHandler.OperationFailure("create objectID", err)
	}
	return docID, nil
}
