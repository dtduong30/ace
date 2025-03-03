package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StringToObjectID converts a string to MongoDB ObjectID
func StringToObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid ID format: %v", err)
	}
	return oid, nil
}
