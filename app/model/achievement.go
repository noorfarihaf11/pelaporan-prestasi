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

type CreateAchievement struct {
    AchievementType string                 `json:"achievement_type" bson:"achievement_type"`
    Title           string                 `json:"title" bson:"title"`
    Description     string                 `json:"description" bson:"description"`
    Details         bson.M                 `json:"details" bson:"details"`
    Tags            []string               `json:"tags" bson:"tags"`
}

type VerifyAchievement struct {
    Points      int       `json:"points" bson:"points"`
    VerifiedBy  string    `json:"verified_by" bson:"verified_by"`
    VerifiedAt  time.Time `json:"verified_at" bson:"verified_at"`
}

type AchievementResponse struct {
    ID              primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
    StudentID       string                 `json:"student_id" bson:"student_id"`
    AdvisorID       string                 `json:"advisor_id" bson:"advisor_id"`
    AchievementType string                 `json:"achievement_type" bson:"achievement_type"`
    Title           string                 `json:"title" bson:"title"`
    Description     string                 `json:"description" bson:"description"`
    Details         bson.M                 `json:"details" bson:"details"`
    Attachments     []Attachment           `json:"attachments" bson:"attachments"`
    Tags            []string               `json:"tags" bson:"tags"`
    Points          int                    `json:"points" bson:"points"`
    Status          string                 `json:"status" bson:"status"`
    VerifiedAt      *time.Time              `json:"verified_at"`
    VerifiedBy      *string                 `json:"verified_by"`
    RejectionNote   *string                 `json:"rejection_note"`
    CreatedAt       time.Time              `json:"created_at" bson:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at" bson:"updated_at"`
}

type Attachment struct {
    FileName   string    `json:"file_name" bson:"file_name"`
    FileUrl    string    `json:"file_url" bson:"file_url"`
    FileType   string    `json:"file_type" bson:"file_type"`
    UploadedAt time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

type AchievementHistory struct {
	Status        string     `json:"status"`
	UpdatedBy     string     `json:"updated_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	RejectionNote *string    `json:"rejection_note,omitempty"`
}

