package repository

import (
	"database/sql"
	"fmt"
	"pelaporan-prestasi/app/model"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func CreateAchievementReference(db *sql.DB, studentID uuid.UUID, mongoID string) error {
    query := `
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, 'draft', $4, $5)
    `

    _, err := db.Exec(query, uuid.New(), studentID, mongoID, time.Now(), time.Now())
    if err != nil {
        return fmt.Errorf("postgres_error: %v", err)
    }

    return nil
}

func UpdateAchievementReference(db *sql.DB, mongoID string, status string) error {
    query := `
        UPDATE achievement_references
        SET status = $1, updated_at = NOW()
        WHERE mongo_achievement_id = $2
        RETURNING id;
    `

    var refID string
    err := db.QueryRow(query, status, mongoID).Scan(&refID)

    if err == sql.ErrNoRows {
        return fmt.Errorf("reference_not_found")
    }

    if err != nil {
        return fmt.Errorf("db_error: %v", err)
    }

    return nil
}
func RejectAchievementReference(db *sql.DB, mongoID string, rejectionNote string) error {
    query := `
        UPDATE achievement_references
        SET status = 'rejected',
            rejection_note = $1,
            updated_at = NOW()
        WHERE mongo_achievement_id = $2
        RETURNING id;
    `

    var refID string
    err := db.QueryRow(query, rejectionNote, mongoID).Scan(&refID)

    if err == sql.ErrNoRows {
        return fmt.Errorf("reference_not_found")
    }

    if err != nil {
        return fmt.Errorf("db_error: %v", err)
    }

    return nil
}

func GetAchievementReferenceByMongoID(db *sql.DB, mongoID string) (*model.AchievementReference, error) {
	query := `
		SELECT id, verified_at, verified_by, rejection_note
		FROM achievement_references
		WHERE mongo_achievement_id = $1
	`

	ref := model.AchievementReference{}
	err := db.QueryRow(query, mongoID).Scan(
		&ref.ID,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
	)

	if err == sql.ErrNoRows {
		return nil, nil 
	}

	if err != nil {
		return nil, err
	}

	return &ref, nil
}
func GetStudentIDByUserID(db *sql.DB, userID string) (string, error) {
    query := `SELECT id FROM students WHERE user_id = $1`

    var id string
    err := db.QueryRow(query, userID).Scan(&id)
    if err == sql.ErrNoRows {
        return "", nil
    }
    if err != nil {
        return "", err
    }
    return id, nil
}
func GetStudentsByAdvisor(db *sql.DB, advisorID string) ([]string, error) {
    query := `SELECT id FROM students WHERE advisor_id = $1`

    // convert advisorID → UUID
    advisorUUID, err := uuid.Parse(advisorID)
    if err != nil {
        fmt.Println("DEBUG invalid advisorID:", advisorID)
        return []string{}, nil
    }

    rows, err := db.Query(query, advisorUUID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var ids []string
    for rows.Next() {
        var studentUUID uuid.UUID  // <-- FIX: pakai uuid.UUID
        if err := rows.Scan(&studentUUID); err != nil {
            return nil, err
        }

        ids = append(ids, studentUUID.String())
    }

    fmt.Println("DEBUG students found for lecturer:", advisorID, "=>", ids)

    return ids, nil
}

func GetAchievementReferencesByStudentIDsPaginated(
    db *sql.DB,
    studentIDs []string,
    limit int,
    offset int,
) ([]model.AchievementReference, int, error) {

    if len(studentIDs) == 0 {
        return []model.AchievementReference{}, 0, nil
    }

    query := `
        SELECT id, student_id, mongo_achievement_id, verified_at, verified_by, rejection_note
        FROM achievement_references
        WHERE student_id = ANY($1)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

    rows, err := db.Query(query, pq.Array(studentIDs), limit, offset)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var list []model.AchievementReference
    for rows.Next() {
        var r model.AchievementReference
        err := rows.Scan(&r.ID, &r.StudentID, &r.MongoAchievementID, &r.VerifiedAt, &r.VerifiedBy, &r.RejectionNote)
        if err != nil {
            return nil, 0, err
        }
        list = append(list, r)
    }

    // Count total
    countQuery := `
        SELECT COUNT(*)
        FROM achievement_references
        WHERE student_id = ANY($1)
    `
    var total int
    err = db.QueryRow(countQuery, pq.Array(studentIDs)).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    return list, total, nil
}
func GetLecturerIDByUserID(db *sql.DB, userID string) (string, error) {
    var lecturerID string

    err := db.QueryRow(`
        SELECT id 
        FROM lecturers 
        WHERE user_id = $1
    `, userID).Scan(&lecturerID)

    if err == sql.ErrNoRows {
        return "", nil
    }
    if err != nil {
        return "", err
    }

    return lecturerID, nil
}
