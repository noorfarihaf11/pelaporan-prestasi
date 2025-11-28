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
    var dto model.UserDTO
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

func UpdateUserService(c *fiber.Ctx, db *sql.DB) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "missing_user_id",
        })
    }

    userUUID, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_user_id_format",
        })
    }

    var dto model.UserDTO
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_request_body",
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

    
    var hashed string
    if dto.Password != "" {
        hashed, err = utils.HashPassword(dto.Password)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_hashing_password",
            })
        }
    }

    
    existingUser, err := repository.GetUserByID(db, id)
    if err != nil || existingUser == nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  "error",
            "message": "user_not_found",
        })
    }

    
    if dto.FullName != "" {
        existingUser.FullName = dto.FullName
    }
    if dto.Username != "" {
        existingUser.Username = dto.Username
    }
    if dto.Email != "" {
        existingUser.Email = dto.Email
    }
    if hashed != "" {
        existingUser.PasswordHash = hashed
    }
    if dto.RoleID != "" {
        roleUUID, err := uuid.Parse(dto.RoleID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "status": "error",
                "message": "invalid_role_id_format",
            })
        }
        existingUser.RoleID = roleUUID
    }

    
    updatedUser, err := repository.UpdateUserTx(tx, userUUID, existingUser)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_user",
        })
    }

    var studentResult *model.Student = nil
    var lecturerResult *model.Lecturer = nil

    
    if dto.StudentProfile != nil {
        studentDTO := dto.StudentProfile

        var advisorUUID *uuid.UUID = nil
        if studentDTO.AdvisorID != nil && *studentDTO.AdvisorID != "" {
            parsed, err := uuid.Parse(*studentDTO.AdvisorID)
            if err != nil {
                return c.Status(400).JSON(fiber.Map{
                    "status":  "error",
                    "message": "invalid_advisor_id_format",
                })
            }
            advisorUUID = &parsed
        }

        s := &model.Student{
            UserID:       updatedUser.ID,
            StudentID:    studentDTO.StudentID,
            ProgramStudy: studentDTO.ProgramStudy,
            AcademicYear: studentDTO.AcademicYear,
            AdvisorID:    advisorUUID,
        }

        err = repository.UpdateStudentTx(tx, s)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_update_student",
            })
        }

        studentResult = s
    }

    if dto.LecturerProfile != nil {
        lecturerDTO := dto.LecturerProfile

        l := &model.Lecturer{
            UserID:     updatedUser.ID,
            LecturerID: lecturerDTO.LecturerID,
            Department: lecturerDTO.Department,
        }

        err = repository.UpdateLecturerTx(tx, l)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_update_lecturer",
            })
        }

        lecturerResult = l
    }

    err = tx.Commit()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_commit_transaction",
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "Berhasil memperbarui user!",
        "data": fiber.Map{
            "user":             updatedUser,
            "student_profile":  studentResult,
            "lecturer_profile": lecturerResult,
        },
    })
}

func DeleteUserService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "missing_user_id",
		})
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid_user_id_format",
		})
	}

	existing, err := repository.GetUserByID(db, id)
	if err != nil || existing == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "user_not_found",
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

	err = repository.DeleteStudentTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_delete_student_profile",
		})
	}

	err = repository.DeleteLecturerTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_delete_lecturer_profile",
		})
	}

	err = repository.DeleteUserTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_delete_user",
		})
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_commit_transaction",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user_deleted_successfully",
	})
}

func UpdateUserRoleService(c *fiber.Ctx, db *sql.DB) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "missing_user_id",
        })
    }

    userUUID, err := uuid.Parse(id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid_user_id_format",
        })
    }

    var payload struct {
        RoleID string `json:"role_id"`
    }

    if err := c.BodyParser(&payload); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid_request_body",
        })
    }

    roleUUID, err := uuid.Parse(payload.RoleID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status": "error",
            "message": "invalid_role_id_format",
        })
    }

    err = repository.UpdateUserRole(db, userUUID, roleUUID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status": "error",
            "message": "failed_update_user_role",
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status": "success",
        "message": "Role user berhasil diperbarui",
    })
}
