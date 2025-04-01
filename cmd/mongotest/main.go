package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Try with no auth
	fmt.Println("Attempting to connect to MongoDB without authentication...")
	err := connectToMongoDB("mongodb://localhost:27017", "battleship", "", "")
	if err != nil {
		fmt.Printf("Failed to connect without auth: %v\n", err)
	} else {
		fmt.Println("Successfully connected without authentication")
	}

	// Try with auth
	fmt.Println("\nAttempting to connect to MongoDB with authentication...")
	err = connectToMongoDB("mongodb://localhost:27017", "battleship", "root", "battleship")
	if err != nil {
		fmt.Printf("Failed to connect with auth: %v\n", err)
	} else {
		fmt.Println("Successfully connected with authentication")
	}
}

func connectToMongoDB(url, dbName, username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uri := url
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", username, password, "localhost:27017", dbName)
	}

	fmt.Printf("Connecting with URI: %s\n", uri)
	
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Try to perform a simple find operation
	collection := client.Database(dbName).Collection("players")
	_, err = collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("error fetching players: %v", err)
	}

	return nil
}