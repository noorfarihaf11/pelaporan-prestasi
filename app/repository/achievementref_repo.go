package repository

import (
    "database/sql"
    "time"

    "github.com/google/uuid"
)

func CreateAchievementReference(db *sql.DB, studentID uuid.UUID, mongoID string) error {
    query := `
        INSERT INTO achievement_references 
        (id, student_id, mongo_achievement_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, 'draft', $4, $5)
    `

    _, err := db.Exec(
        query,
        uuid.New(),       // id
        studentID,
        mongoID,          // hex string MongoDB _id
        time.Now(),
        time.Now(),
    )

    return err
}

func UpdateAchievementReference(db *sql.DB, mongoID string, status string) error {
    query := `
        UPDATE achievement_references
        SET status = $1, updated_at = NOW()
        WHERE mongo_achievement_id = $2
    `
    _, err := db.Exec(query, status, mongoID)
    return err
}
