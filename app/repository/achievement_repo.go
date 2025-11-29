package repository

import (
    "context"
    "fmt"
    "time"
    "pelaporan-prestasi/app/model"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

func CreateAchievement(db *mongo.Database, ach *model.Achievement) (*model.Achievement, error) {
    collection := db.Collection("achievements")

    ach.CreatedAt = time.Now()
    ach.UpdatedAt = time.Now()

    result, err := collection.InsertOne(context.Background(), ach)
    if err != nil {
        return nil, err
    }

    oid, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
    }

    ach.ID = oid
    return ach, nil
}

func GetAllAchievements(db *mongo.Database) ([]model.Achievement, error) {
    collection := db.Collection("achievements")

    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var list []model.Achievement
    for cursor.Next(context.Background()) {
        var a model.Achievement
        if err := cursor.Decode(&a); err != nil {
            return nil, err
        }
        list = append(list, a)
    }

    return list, nil
}

func GetAchievementByID(db *mongo.Database, id string) (*model.Achievement, error) {
    collection := db.Collection("achievements")

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var ach model.Achievement
    err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&ach)

    if err == mongo.ErrNoDocuments {
        return nil, nil
    }

    if err != nil {
        return nil, err
    }

    return &ach, nil
}
