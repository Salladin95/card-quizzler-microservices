package user

import (
	"github.com/Salladin95/goErrorHandler"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// toObjectId converts a hexadecimal string representation of MongoDB ObjectID to primitive.ObjectID.
// It returns the converted ObjectID and an error.
// If the conversion fails, an error with context is returned.
func toObjectId(id string) (primitive.ObjectID, error) {
	// Convert the hexadecimal string to primitive.ObjectID
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return docID, goErrorHandler.OperationFailure("create objectID", err)
	}
	return docID, nil
}
