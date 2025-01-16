package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

type Wishlist struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`                  // e.g., "Read", "Buy", "Learn", "Do"
	Description string    `json:"description,omitempty"` // Optional extra details
	Priority    int       `json:"priority,omitempty"`    // 1 (low) - 5 (high) priority
	Status      string    `json:"status"`                // "pending", "in-progress", "completed"
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

var client *mongo.Client

func main() {
	// Kết nối MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb+srv://deadpool:vippro@cluster0.lir7snr.mongodb.net/test?retryWrites=true&w=majority" // fallback
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, _ = mongo.Connect(context.Background(), clientOptions)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())

	r := gin.Default()

	// CRUD endpoints
	r.GET("/wishlist", getWishlist)
	r.POST("/wishlist", createWishlist)
	r.PUT("/wishlist/:id", updateWishlist)
	r.DELETE("/wishlist/:id", deleteWishlist)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func getWishlist(c *gin.Context) {
	collection := client.Database("ace").Collection("wishlists")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var wishlist []Wishlist
	if err = cursor.All(context.Background(), &wishlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

func createWishlist(c *gin.Context) {
	var wishlist Wishlist
	if err := c.ShouldBindJSON(&wishlist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	wishlist.CreatedAt = time.Now()
	wishlist.UpdatedAt = time.Now()

	collection := client.Database("ace").Collection("wishlists")
	result, err := collection.InsertOne(context.Background(), wishlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed: " + err.Error()})
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		wishlist.ID = oid.Hex()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert ID"})
		return
	}

	c.JSON(http.StatusOK, wishlist)
}

func updateWishlist(c *gin.Context) {
	id := c.Param("id")
	var wishlist Wishlist
	err := c.BindJSON(&wishlist)
	if err != nil {
		return
	}
	collection := client.Database("ace").Collection("wishlists")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": wishlist})
	if err != nil {
		return
	}
	c.JSON(200, gin.H{"message": "Wishlist updated"})
}

func deleteWishlist(c *gin.Context) {
	id := c.Param("id")
	collection := client.Database("ace").Collection("wishlist")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return
	}
	c.JSON(200, gin.H{"message": "Wishlist deleted"})
}
