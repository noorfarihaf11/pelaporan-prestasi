package model

import (
	"github.com/google/uuid"
	"time"
)

type Lecturer struct {
	ID           	uuid.UUID	`json:"id"`                       
	UserID       	uuid.UUID	`json:"user_id"`       
	LecturerID    	string		`json:"lecturer_id"` 
	Department    	string	 	`json:"department"` 
	CreatedAt   	time.Time 	`json:"created_at"`
}
type StudentResponse struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
}