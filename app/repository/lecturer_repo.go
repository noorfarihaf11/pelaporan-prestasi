package repository

import (
	"database/sql"
	"pelaporan-prestasi/app/model"

	"github.com/google/uuid"
)

func CreateLecturerTx(tx *sql.Tx, lecturer *model.Lecturer) error {
    query := `
        INSERT INTO lecturers (user_id, lecturer_id, department)
        VALUES ($1, $2, $3)
        RETURNING id
    `
    err := tx.QueryRow(query,
        lecturer.UserID,
        lecturer.LecturerID,
        lecturer.Department,
    ).Scan(&lecturer.ID)

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
func GetAdviseesByLecturerID(db *sql.DB, lecturerID string) ([]model.StudentResponse, error) {
	advisorUUID, err := uuid.Parse(lecturerID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT s.id, u.full_name, s.student_id, s.program_study
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
		ORDER BY u.full_name
	`

	rows, err := db.Query(query, advisorUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advisees []model.StudentResponse
	for rows.Next() {
		var s model.StudentResponse
		if err := rows.Scan(&s.ID, &s.FullName, &s.StudentID, &s.ProgramStudy); err != nil {
			return nil, err
		}
		advisees = append(advisees, s)
	}

	return advisees, nil
}
