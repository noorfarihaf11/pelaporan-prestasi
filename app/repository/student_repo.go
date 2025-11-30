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
        RETURNING id
    `

	return tx.QueryRow(
		query,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.CreatedAt,
	).Scan(&student.ID)
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

func GetAllStudent(db *sql.DB) ([]model.Student, error) {
	rows, err := db.Query(`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
    FROM students`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentList []model.Student
	for rows.Next() {
		var s model.Student
		err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		if err != nil {
			return nil, err
		}
		studentList = append(studentList, s)
	}

	return studentList, nil
}

func GetStudentByID(db *sql.DB, id string) (*model.Student, error) {
	row := db.QueryRow(`SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at
    FROM students WHERE id=$1`, id)

	var s model.Student
	err := row.Scan(
	&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}
