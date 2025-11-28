package repository

import (
	"database/sql"
	"pelaporan-prestasi/app/model"

	"github.com/google/uuid"
)

func CreateStudentTx(tx *sql.Tx, student *model.Student) error {
    query := `
        INSERT INTO students (user_id, student_id, program_study, academic_year, advisor_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := tx.Exec(
        query,
        student.UserID,
        student.StudentID,
        student.ProgramStudy,
        student.AcademicYear,
        student.AdvisorID,
        student.CreatedAt,
    )

    return err
}

func UpdateStudentTx(tx *sql.Tx, s *model.Student) error {
    _, err := tx.Exec(`
        UPDATE students
        SET student_id=$2, program_study=$3, academic_year=$4, advisor_id=$5
        WHERE user_id=$1
    `,
        
        s.UserID,
        s.StudentID,
        s.ProgramStudy,
        s.AcademicYear,
        s.AdvisorID,
    )
    return err
}

func DeleteStudentTx(tx *sql.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM students WHERE user_id=$1`, userID)
	return err
}
