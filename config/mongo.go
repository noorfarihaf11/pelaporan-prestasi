package config

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Database {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Mongo connect error:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Mongo ping error:", err)
    }

    log.Println("MongoDB connected successfully!")

    return client.Database("pelaporan-prestasi")
}
