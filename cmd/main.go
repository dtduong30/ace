package main

import (
	"ace/internal/handlers"
	"ace/internal/repository"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	log.Println("URI:", uri)
	if uri == "" {
		uri = "mongodb://localhost:27017/"
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {

		}
	}(client, context.Background())
	log.Println("🚀 MongoDB connected")

	// Init repository và handler
	wishlistRepo := repository.NewWishlistRepository(client)
	wishlistHandler := handlers.NewWishlistHandler(wishlistRepo)

	// Setup Gin router
	r := gin.Default()
	r.GET("/wishlist", wishlistHandler.GetAllWishlist)
	r.POST("/wishlist", wishlistHandler.CreateWishlist)
	r.PUT("/wishlist/:id", wishlistHandler.UpdateWishlist)
	r.GET("/wishlist/:id", wishlistHandler.GetOneWishlist)
	r.DELETE("/wishlist/:id", wishlistHandler.DeleteWishlist)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
