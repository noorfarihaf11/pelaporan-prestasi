package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pelaporan-prestasi/app/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func AddAttachmentsToAchievement(mongoDB *mongo.Database, achievementID string, attachments []model.Attachment) error {
	collection := mongoDB.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"attachments": bson.M{
				"$each": attachments,
			},
		},
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	return err
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
			a.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
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
		ach.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
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

func VerifyAchievement(mongoDB *mongo.Database, achievementID string, points int, lecturerID string) error {
	collection := mongoDB.Collection("achievements")

	objID, err := primitive.ObjectIDFromHex(achievementID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":      "verified",
			"points":      points,
			"verified_by": lecturerID,
			"verified_at": time.Now(),
			"updated_at":  time.Now(),
		},
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	return err
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
			ach.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
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
			updatedBy = *ach.VerifiedBy
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
func GetAchievementsForStudent(
	mongoDB *mongo.Database,
	sqlDB *sql.DB,
	studentID string,
) ([]model.AchievementResponse, error) {

	collection := mongoDB.Collection("achievements")

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

		// merge reference from Postgres
		ref, err := GetAchievementReferenceByMongoID(sqlDB, ach.ID.Hex())
		if err != nil {
			return nil, err
		}
		if ref != nil {
			ach.VerifiedAt = ref.VerifiedAt
			ach.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
			ach.RejectionNote = ref.RejectionNote
		}

		// advisor_id
		studentUUID, err := uuid.Parse(studentID)
		if err == nil {
			sqlDB.QueryRow("SELECT advisor_id FROM students WHERE id=$1", studentUUID).Scan(&ach.AdvisorID)
		}

		achievements = append(achievements, ach)
	}

	return achievements, nil
}
func GetLecturerAchievementsPaginated(
	sqlDB *sql.DB,
	mongoDB *mongo.Database,
	lecturerID string,
	limit int,
	offset int,
) ([]model.AchievementResponse, int, error) {

	// Ambil student ID
	studentIDs, err := GetStudentsByAdvisor(sqlDB, lecturerID)
	if err != nil {
		return nil, 0, err
	}
	if len(studentIDs) == 0 {
		return []model.AchievementResponse{}, 0, nil
	}

	// Query paginated achievement references
	query := `
        SELECT student_id, mongo_achievement_id, verified_at, verified_by, rejection_note
        FROM achievement_references
        WHERE student_id = ANY($1)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
	rows, err := sqlDB.Query(query, pq.Array(studentIDs), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	mongoColl := mongoDB.Collection("achievements")
	list := []model.AchievementResponse{}

	for rows.Next() {
		var ref model.AchievementReference
		if err := rows.Scan(&ref.StudentID, &ref.MongoAchievementID, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote); err != nil {
			return nil, 0, err
		}

		oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		var ach model.AchievementResponse
		mongoColl.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&ach)

		ach.StudentID = ref.StudentID.String()
		ach.VerifiedAt = ref.VerifiedAt
		ach.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
		ach.RejectionNote = ref.RejectionNote
		ach.AdvisorID = lecturerID

		list = append(list, ach)
	}

	// Count total entries
	countQuery := `
        SELECT COUNT(*)
        FROM achievement_references
        WHERE student_id = ANY($1)
    `
	var total int
	err = sqlDB.QueryRow(countQuery, pq.Array(studentIDs)).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func UUIDPtrToStringPtr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}
func GetRoleNameByID(db *sql.DB, roleID string) (string, error) {
	var roleName string

	// Convert string → UUID
	uuidVal, err := uuid.Parse(roleID)
	if err != nil {
		return "", fmt.Errorf("invalid_role_uuid")
	}

	err = db.QueryRow("SELECT name FROM roles WHERE id = $1", uuidVal).Scan(&roleName)
	if err != nil {
		return "", err
	}

	return roleName, nil
}
func GetAdminAchievementsPaginated(
	sqlDB *sql.DB,
	mongoDB *mongo.Database,
	limit int,
	offset int,
	sort string,
	order string,
	status string,
	studentName string,
) ([]model.AchievementResponse, int, error) {

	filters := []string{"1 = 1"}
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		filters = append(filters, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, status)
		argIndex++
	}

	if studentName != "" {
		filters = append(filters, fmt.Sprintf("student_name ILIKE $%d", argIndex))
		args = append(args, "%"+studentName+"%")
		argIndex++
	}

	filterQuery := strings.Join(filters, " AND ")

	query := fmt.Sprintf(`
        SELECT student_id, mongo_achievement_id, verified_at, verified_by, rejection_note
        FROM achievement_references
        WHERE %s
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d
    `, filterQuery, sort, order, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := sqlDB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	mongoColl := mongoDB.Collection("achievements")
	list := []model.AchievementResponse{}

	for rows.Next() {
		var ref model.AchievementReference
		rows.Scan(&ref.StudentID, &ref.MongoAchievementID, &ref.VerifiedAt, &ref.VerifiedBy, &ref.RejectionNote)

		oid, _ := primitive.ObjectIDFromHex(ref.MongoAchievementID)
		var ach model.AchievementResponse
		mongoColl.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&ach)

		ach.StudentID = ref.StudentID.String()
		ach.VerifiedAt = ref.VerifiedAt
		ach.VerifiedBy = UUIDPtrToStringPtr(ref.VerifiedBy)
		ach.RejectionNote = ref.RejectionNote

		list = append(list, ach)
	}

	// total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM achievement_references WHERE %s`, filterQuery)
	var total int
	sqlDB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)

	return list, total, nil
}
