package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
    ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
    StudentID       string                 `json:"student_id" bson:"student_id"`
    AchievementType string                 `json:"achievement_type" bson:"achievement_type"`
    Title           string                 `json:"title" bson:"title"`
    Description     string                 `json:"description" bson:"description"`
    Details         bson.M                 `json:"details" bson:"details"`
    Attachments     []Attachment           `json:"attachments" bson:"attachments"`
    Tags            []string               `json:"tags" bson:"tags"`
    Points          int                    `json:"points" bson:"points"`
    Status          string                 `json:"status" bson:"status"`
    CreatedAt       time.Time              `json:"created_at" bson:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at" bson:"updated_at"`
}


type Attachment struct {
    FileName   string    `json:"file_name" bson:"file_name"`
    FileUrl    string    `json:"file_url" bson:"file_url"`
    FileType   string    `json:"file_type" bson:"file_type"`
    UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
}
