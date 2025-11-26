package repository

import (
    "database/sql"
    "pelaporan-prestasi/app/model"
)

func CreateLecturerTx(tx *sql.Tx, lecturer *model.Lecturer) error {
    query := `
        INSERT INTO lecturers (user_id, lecturer_id, department, created_at)
        VALUES ($1, $2, $3, $4)
    `

    _, err := tx.Exec(
        query,
        lecturer.UserID,
        lecturer.LecturerID,
        lecturer.Department,
        lecturer.CreatedAt,
    )

    return err
}

func UpdateLecturerTx(tx *sql.Tx, l *model.Lecturer) error {
    _, err := tx.Exec(`
        UPDATE lecturers
        SET lecturer_id=$2, department=$3
        WHERE user_id=$1
    `,
        
        l.UserID,
        l.LecturerID,
        l.Department,
    )
    return err
}
