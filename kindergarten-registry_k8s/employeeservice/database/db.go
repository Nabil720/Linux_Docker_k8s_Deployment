package database

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() error {
    // Environment variables থেকে MongoDB URI নিবে
    // (Vault থেকে main.go-তে already set করা হয়েছে)
    connectionString := os.Getenv("MONGODB_URI")
    if connectionString == "" {
        log.Fatal("MONGODB_URI environment variable is not set")
    }

    databaseName := os.Getenv("DATABASE_NAME")
    if databaseName == "" {
        databaseName = "kindergarten"
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(connectionString)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    // Ping the database to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    Client = client
    Database = client.Database(databaseName)
    log.Println("Connected to MongoDB successfully!")
    return nil
}

func GetCollection(collectionName string) *mongo.Collection {
    return Database.Collection(collectionName)
}

func Disconnect() {
    if Client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        Client.Disconnect(ctx)
        log.Println("MongoDB connection closed")
    }
}