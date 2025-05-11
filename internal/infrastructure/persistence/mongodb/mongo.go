package persistence

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance        *mongo.Client
	clientOnce            sync.Once
	MongoConnectionString = "MONGODB_URI"
)

// ConnectMongoDB khởi tạo và trả về kết nối MongoDB
func ConnectMongoDB(uri string) *mongo.Client {
	clientOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Kiểm tra kết nối bằng Ping
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		clientInstance = client
		fmt.Println("✅ Connected to MongoDB!")
	})

	return clientInstance
}


// GetMongoClient trả về client MongoDB đã kết nối
func GetMongoClient() *mongo.Client {
	if clientInstance == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	return clientInstance
}
