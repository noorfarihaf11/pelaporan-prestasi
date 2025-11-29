package model

import (
    "time"

    "github.com/google/uuid"
)

type AchievementReference struct {
    ID                 uuid.UUID
    StudentID          uuid.UUID
    MongoAchievementID string
    Status             string
    SubmittedAt        *time.Time
    VerifiedAt         *time.Time
    VerifiedBy         *uuid.UUID
    RejectionNote      *string
    CreatedAt          time.Time
    UpdatedAt          time.Time
}