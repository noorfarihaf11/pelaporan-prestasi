package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pelaporan-prestasi/app/model"
	"time"

	"github.com/google/uuid"
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

func GetAllAchievements(db *mongo.Database, sqlDB *sql.DB) ([]model.AchievementResponse, error) {
	collection := db.Collection("achievements")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var list []model.AchievementResponse
	for cursor.Next(context.Background()) {
		var a model.AchievementResponse
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}

		var advisorID string
		studentUUID, err := uuid.Parse(a.StudentID) // a.StudentID harus UUID
		if err != nil {
			// Kalau StudentID tidak valid UUID, skip atau set advisorID kosong
			a.AdvisorID = ""
		} else {
			err = sqlDB.QueryRow("SELECT advisor_id FROM students WHERE id=$1", studentUUID).Scan(&advisorID)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}
			a.AdvisorID = advisorID
		}

		// Ambil reference dari PostgreSQL
		ref, err := GetAchievementReferenceByMongoID(sqlDB, a.ID.Hex())
		if err != nil {
			return nil, err
		}
		if ref != nil {
			a.VerifiedAt = ref.VerifiedAt
			a.VerifiedBy = ref.VerifiedBy
			a.RejectionNote = ref.RejectionNote
		}

		list = append(list, a)
	}

	return list, nil
}

func GetAchievementByID(db *mongo.Database, sqlDB *sql.DB, id string) (*model.AchievementResponse, error) {
	collection := db.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var ach model.AchievementResponse
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&ach)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Ambil reference
	ref, err := GetAchievementReferenceByMongoID(sqlDB, id)
	if err != nil {
		return nil, err
	}
	if ref != nil {
		ach.VerifiedAt = ref.VerifiedAt
		ach.VerifiedBy = ref.VerifiedBy
		ach.RejectionNote = ref.RejectionNote
	}

	var advisorID string
	studentUUID, err := uuid.Parse(ach.StudentID)
	if err == nil {
		_ = sqlDB.QueryRow("SELECT advisor_id FROM students WHERE id=$1", studentUUID).Scan(&advisorID)
	}
	ach.AdvisorID = advisorID

	return &ach, nil
}

// Helper function untuk ambil advisor_id
func GetAdvisorIDByStudentID(db *sql.DB, studentID string) (*string, error) {
	var advisorID *string
	query := `SELECT advisor_id FROM students WHERE id = $1`
	err := db.QueryRow(query, studentID).Scan(&advisorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return advisorID, nil
}

func UpdateAchievement(db *mongo.Database, id string, update bson.M) (*model.Achievement, error) {
	collection := db.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	var updated model.Achievement
	collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&updated)

	return &updated, nil
}

func SoftDeleteAchievement(db *mongo.Database, id string) error {
	collection := db.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"status": "deleted",
			},
		},
	)
	return err
}

func UpdateAchievementStatus(db *mongo.Database, id string, status string) error {
	collection := db.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid_id: %v", err)
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("mongo_update_error: %v", err)
	}

	return nil
}

func GetAchievementsByStudentID(db *mongo.Database, sqlDB *sql.DB, studentID string) ([]model.AchievementResponse, error) {
	collection := db.Collection("achievements")

	cursor, err := collection.Find(context.Background(), bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var achievements []model.AchievementResponse
	for cursor.Next(context.Background()) {
		var ach model.AchievementResponse
		if err := cursor.Decode(&ach); err != nil {
			return nil, err
		}

		ref, err := GetAchievementReferenceByMongoID(sqlDB, ach.ID.Hex())
		if err != nil {
			return nil, err
		}
		if ref != nil {
			ach.VerifiedAt = ref.VerifiedAt
			ach.VerifiedBy = ref.VerifiedBy
			ach.RejectionNote = ref.RejectionNote
		}

		var advisorID string
		studentUUID, err := uuid.Parse(studentID)
		if err == nil {
			_ = sqlDB.QueryRow("SELECT advisor_id FROM students WHERE id=$1", studentUUID).Scan(&advisorID)
		}
		ach.AdvisorID = advisorID

		achievements = append(achievements, ach)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return achievements, nil
}
func GetAchievementHistory(db *sql.DB, mongoDB *mongo.Database, mongoID string) ([]model.AchievementHistory, error) {
	// Ambil status dari MongoDB
	ach, err := GetAchievementByID(mongoDB, db, mongoID)
	if err != nil {
		return nil, err
	}
	if ach == nil {
		return nil, nil
	}

	history := []model.AchievementHistory{
		{
			Status:    ach.Status,
			UpdatedBy: ach.StudentID,
			UpdatedAt: ach.UpdatedAt,
		},
	}

	// Ambil verified info jika ada
	if ach.VerifiedAt != nil {
		var updatedBy string
		if ach.VerifiedBy != nil {
			updatedBy = ach.VerifiedBy.String()
		}
		history = append(history, model.AchievementHistory{
			Status:        "verified",
			UpdatedBy:     updatedBy,
			UpdatedAt:     *ach.VerifiedAt,
			RejectionNote: ach.RejectionNote,
		})
	}

	return history, nil
}
