package repository

import (
	"ace/internal/models"
	"ace/internal/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type WishlistRepository struct {
	Collection *mongo.Collection
}

// NewWishlistRepository creates a new repository instance
func NewWishlistRepository(client *mongo.Client) *WishlistRepository {
	return &WishlistRepository{
		Collection: client.Database("ace").Collection("wishlists"),
	}
}

// Create inserts a new wishlist item
func (r *WishlistRepository) Create(ctx context.Context, wishlist *models.Wishlist) (string, error) {
	wishlist.CreatedAt = time.Now()
	wishlist.UpdatedAt = time.Now()

	result, err := r.Collection.InsertOne(ctx, wishlist)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", nil
}

// Update updates an existing wishlist item
func (r *WishlistRepository) Update(ctx context.Context, id string, wishlist *models.Wishlist) error {
	oid, err := utils.StringToObjectID(id)
	if err != nil {
		return err
	}
	wishlist.UpdatedAt = time.Now()
	result, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"title":       wishlist.Title,
			"type":        wishlist.Type,
			"description": wishlist.Description,
			"priority":    wishlist.Priority,
			"status":      wishlist.Status,
			"updated_at":  wishlist.UpdatedAt,
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with id: %s", id)
	}
	return nil
}

func (r *WishlistRepository) GetAll(ctx context.Context, page, limit int) ([]models.Wishlist, int64, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := r.Collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var wishlists []models.Wishlist
	if err = cursor.All(ctx, &wishlists); err != nil {
		return nil, 0, err
	}

	total, err := r.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return wishlists, total, nil
}

func (r *WishlistRepository) GetOne(ctx context.Context, id string) (*models.Wishlist, error) {
	oid, err := utils.StringToObjectID(id)
	if err != nil {
		return nil, err
	}
	var wishlist models.Wishlist
	err = r.Collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&wishlist)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no document found with id: %s", id)
	}
	if err != nil {
		return nil, err
	}
	return &wishlist, nil
}

func (r *WishlistRepository) Delete(ctx context.Context, id string) error {
	oid, err := utils.StringToObjectID(id)
	if err != nil {
		return err
	}
	result, err := r.Collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with id: %s", id)
	}
	return nil
}
