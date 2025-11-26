package service

import (
	"database/sql"
	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAllUserService(c *fiber.Ctx, db *sql.DB) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Token",
		})
	}

	userList, err := repository.GetAllUser(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetching users",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mendapatkan semua data user",
		"data": fiber.Map{
			"users": userList,
		},
	})
}
func GetUserByIDService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing user ID",
		})
	}

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	_, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	user, err := repository.GetUserByID(db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetching user",
		})
	}

	if user == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User fetched successfully",
		"data": fiber.Map{
			"user": user,
		},
	})
}

func CreateUserService(c *fiber.Ctx, db *sql.DB) error {
    var dto model.CreateUserDTO
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_request_body",
        })
    }

    roleUUID, err := uuid.Parse(dto.RoleID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_role_id_format",
        })
    }

    tx, err := db.Begin()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_to_start_transaction",
        })
    }

    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    hashed, err := utils.HashPassword(dto.Password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_hashing_password",
        })
    }

    user := &model.User{
        FullName:     dto.FullName,
        Username:     dto.Username,
        Email:        dto.Email,
        PasswordHash: hashed,
        RoleID:       roleUUID,
        IsActive:     true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    user, err = repository.CreateUserTx(tx, user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_creating_user",
        })
    }

    var student_result *model.Student = nil
    var lecturer_result *model.Lecturer = nil

    if dto.StudentProfile != nil {
          var advisorUUID *uuid.UUID = nil

        if dto.StudentProfile.AdvisorID != nil {
            if *dto.StudentProfile.AdvisorID != "" {
                parsed, err := uuid.Parse(*dto.StudentProfile.AdvisorID)
                if err != nil {
                    return c.Status(400).JSON(fiber.Map{
                        "status":  "error",
                        "message": "invalid_advisor_id_format",
                    })
                }
                advisorUUID = &parsed
            }
        }

        student_profile := &model.Student{
            UserID:       user.ID,
            StudentID:    dto.StudentProfile.StudentID,   
            ProgramStudy: dto.StudentProfile.ProgramStudy,
            AcademicYear: dto.StudentProfile.AcademicYear,
            AdvisorID:    advisorUUID,                   
            CreatedAt:    time.Now(),
        }

        err = repository.CreateStudentTx(tx, student_profile)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_creating_student_profile",
            })
        }

        student_result = student_profile
    }

    if dto.LecturerProfile != nil {
        lecturer_profile := &model.Lecturer{
            UserID:     user.ID,
            LecturerID: dto.LecturerProfile.LecturerID,  
            Department: dto.LecturerProfile.Department,
            CreatedAt:  time.Now(),
        }

        err = repository.CreateLecturerTx(tx, lecturer_profile)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_creating_lecturer_profile",
            })
        }
        
        lecturer_result = lecturer_profile
    }

    err = tx.Commit()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_committing_transaction",
        })
    }

    return c.Status(201).JSON(fiber.Map{
        "status":  "success",
        "message": "user_created_successfully",
        "data": fiber.Map{
            "user": user,
            "student_profile": student_result,
            "lecturer_profile": lecturer_result,
        },
    })
}
