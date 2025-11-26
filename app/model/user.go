package model

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserDTO struct {
    FullName  string `json:"full_name"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"password"`
    RoleID    string `json:"role_id"`

    StudentProfile  *StudentProfileDTO  `json:"student_profile,omitempty"`
    LecturerProfile *LecturerProfileDTO `json:"lecturer_profile,omitempty"`
}

type StudentProfileDTO struct {
    StudentID    string `json:"student_id"`
    ProgramStudy string `json:"program_study"`
    AcademicYear string `json:"academic_year"`
    AdvisorID    *string `json:"advisor_id, omitempty"`
}

type LecturerProfileDTO struct {
    LecturerID string `json:"lecturer_id"`
    Department string `json:"department"`
}
