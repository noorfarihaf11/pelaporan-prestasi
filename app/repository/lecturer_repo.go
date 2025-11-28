package repository

import (
	"database/sql"
	"pelaporan-prestasi/app/model"

	"github.com/google/uuid"
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

func DeleteLecturerTx(tx *sql.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM lecturers WHERE user_id=$1`, userID)
	return err
}

func GetAllLecturers(db *sql.DB) ([]model.Lecturer, error) {
	rows, err := db.Query(`SELECT id, user_id, lecturer_id, department, created_at
    FROM lecturers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturerList []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		lecturerList = append(lecturerList, l)
	}

	return lecturerList, nil
}