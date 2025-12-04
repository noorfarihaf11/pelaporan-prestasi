package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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


