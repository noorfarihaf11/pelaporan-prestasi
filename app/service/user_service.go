package service

import (
	"database/sql"
	"errors"
	"fmt"
	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Mengambil seluruh data user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.BaseResponse{data=map[string][]model.User}
// @Failure      401 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users [get]
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

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Mengambil detail user berdasarkan ID
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200 {object} model.BaseResponse{data=map[string]model.User}
// @Failure      400 {object} model.BaseResponse
// @Failure      401 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users/{id} [get]
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

var (
	beginTxFunc          = func(db *sql.DB) (*sql.Tx, error) { return db.Begin() }
	commitTxFunc         = func(tx *sql.Tx) error { return tx.Commit() }
	rollbackTxFunc       = func(tx *sql.Tx) error { return tx.Rollback() }
	hashPasswordFunc     = utils.HashPassword
	createUserTxFunc     = repository.CreateUserTx
	createStudentTxFunc  = repository.CreateStudentTx
	createLecturerTxFunc = repository.CreateLecturerTx
	nowFunc              = time.Now
)
var (
	getUserByIDFunc      = repository.GetUserByID
	updateUserTxFunc    = repository.UpdateUserTx
	updateStudentTxFunc = repository.UpdateStudentTx
	updateLecturerTxFunc = repository.UpdateLecturerTx

	deleteStudentTxFunc = repository.DeleteStudentTx
	deleteLecturerTxFunc = repository.DeleteLecturerTx
	deleteUserTxFunc    = repository.DeleteUserTx
)

// CreateUser godoc
// @Summary      Create new user
// @Description  Membuat user baru (beserta student / lecturer profile jika ada)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      model.UserDTO  true  "Create User Payload"
// @Success      201 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users [post]
func CreateUserService(c *fiber.Ctx, db *sql.DB) error {
	var dto model.UserDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	roleUUID, err := uuid.Parse(dto.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid role id format",
		})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to start transaction",
		})
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	hashed, err := utils.HashPassword(dto.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed hasing password",
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
			"message": "Failed creating user",
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
						"message": "Invalid advisor id format",
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
			fmt.Println("Repository error:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed creating student profile",
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
				"message": "Failed creating lecturer profile",
			})
		}

		lecturer_result = lecturer_profile
	}

	err = commitTxFunc(tx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed commiting transaction",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "User created successfully",
		"data": fiber.Map{
			"user":             user,
			"student_profile":  student_result,
			"lecturer_profile": lecturer_result,
		},
	})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Memperbarui data user beserta student / lecturer profile
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string        true  "User ID"
// @Param        body  body      model.UserDTO true  "Update User Payload"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users/{id} [put]
func UpdateUserService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing user ID",
		})
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user id format",
		})
	}

	var dto model.UserDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to start transactions",
		})
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	var hashed string
	if dto.Password != "" {
		hashed, err = utils.HashPassword(dto.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed hasing password",
			})
		}
	}

	existingUser, err := repository.GetUserByID(db, id)
	if err != nil || existingUser == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
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
				"status":  "error",
				"message": "Invalid role id format",
			})
		}
		existingUser.RoleID = roleUUID
	}

	updatedUser, err := repository.UpdateUserTx(tx, userUUID, existingUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update user",
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
					"message": "Invalid advisor id format",
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
				"message": "Failed update student",
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
				"message": "Failed update lecturer",
			})
		}

		lecturerResult = l
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed commit transaction",
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

// DeleteUser godoc
// @Summary      Delete user
// @Description  Menghapus user beserta relasi student dan lecturer
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users/{id} [delete]
func DeleteUserService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing user id",
		})
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user id format",
		})
	}

	existing, err := repository.GetUserByID(db, id)
	if err != nil || existing == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to start transaction",
		})
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	err = repository.DeleteStudentTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed delete student",
		})
	}

	err = repository.DeleteLecturerTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed deleted lecturer",
		})
	}

	err = repository.DeleteUserTx(tx, userUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed delete user",
		})
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed commit transaction",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User deleted succesfully",
	})
}

// UpdateUserRole godoc
// @Summary      Update user role
// @Description  Mengubah role user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true  "User ID"
// @Param        body  body      object  true  "Role Payload"  example({"role_id":"uuid"})
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/users/{id}/role [patch]
func UpdateUserRoleService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing user id",
		})
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user id format",
		})
	}

	var payload struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	roleUUID, err := uuid.Parse(payload.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid role id format",
		})
	}

	err = repository.UpdateUserRole(db, userUUID, roleUUID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update user profile",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Role user berhasil diperbarui",
	})
}

// test logic, karna fiber tidak bisa dilakukan unit testing
func CreateUser(db *sql.DB, dto model.UserDTO) (*model.User, *model.Student, *model.Lecturer, error) {
	roleUUID, err := uuid.Parse(dto.RoleID)
	if err != nil {
		return nil, nil, nil, errors.New("invalid role id format")
	}

	tx, err := beginTxFunc(db)
	if err != nil {
		return nil, nil, nil, err
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	hashed, err := hashPasswordFunc(dto.Password)
	if err != nil {
		return nil, nil, nil, err
	}

	user := &model.User{
		FullName:     dto.FullName,
		Username:     dto.Username,
		Email:        dto.Email,
		PasswordHash: hashed,
		RoleID:       roleUUID,
		IsActive:     true,
		CreatedAt:    nowFunc(),
		UpdatedAt:    nowFunc(),
	}

	user, err = createUserTxFunc(tx, user)
	if err != nil {
		return nil, nil, nil, err
	}

	var student *model.Student
	var lecturer *model.Lecturer

	if dto.StudentProfile != nil {
		student = &model.Student{
			UserID:       user.ID,
			StudentID:    dto.StudentProfile.StudentID,
			ProgramStudy: dto.StudentProfile.ProgramStudy,
			AcademicYear: dto.StudentProfile.AcademicYear,
			CreatedAt:    nowFunc(),
		}

		err = createStudentTxFunc(tx, student)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if dto.LecturerProfile != nil {
		lecturer = &model.Lecturer{
			UserID:     user.ID,
			LecturerID: dto.LecturerProfile.LecturerID,
			Department: dto.LecturerProfile.Department,
			CreatedAt:  nowFunc(),
		}

		err = createLecturerTxFunc(tx, lecturer)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	err = commitTxFunc(tx)
	if err != nil {
		return nil, nil, nil, err
	}

	return user, student, lecturer, nil
}
func UpdateUser(db *sql.DB, userID string, dto model.UserDTO) (
	*model.User, *model.Student, *model.Lecturer, error,
) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, nil, nil, err
	}

	tx, err := beginTxFunc(db)
	if err != nil {
		return nil, nil, nil, err
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	existingUser, err := getUserByIDFunc(db, userID)
	if err != nil || existingUser == nil {
		return nil, nil, nil, errors.New("user not found")
	}

	if dto.Password != "" {
		existingUser.PasswordHash, err = hashPasswordFunc(dto.Password)
		if err != nil {
			return nil, nil, nil, err
		}
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
	if dto.RoleID != "" {
		existingUser.RoleID, err = uuid.Parse(dto.RoleID)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	updatedUser, err := updateUserTxFunc(tx, userUUID, existingUser)
	if err != nil {
		return nil, nil, nil, err
	}

	var student *model.Student
	var lecturer *model.Lecturer

	if dto.StudentProfile != nil {
		student = &model.Student{
			UserID:       updatedUser.ID,
			StudentID:    dto.StudentProfile.StudentID,
			ProgramStudy: dto.StudentProfile.ProgramStudy,
			AcademicYear: dto.StudentProfile.AcademicYear,
		}

		err = updateStudentTxFunc(tx, student)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if dto.LecturerProfile != nil {
		lecturer = &model.Lecturer{
			UserID:     updatedUser.ID,
			LecturerID: dto.LecturerProfile.LecturerID,
			Department: dto.LecturerProfile.Department,
		}

		err = updateLecturerTxFunc(tx, lecturer)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	err = commitTxFunc(tx)
	if err != nil {
		return nil, nil, nil, err
	}

	return updatedUser, student, lecturer, nil
}

func DeleteUser(db *sql.DB, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	existing, err := getUserByIDFunc(db, userID)
	if err != nil || existing == nil {
		return errors.New("user not found")
	}

	tx, err := beginTxFunc(db)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackTxFunc(tx)
		}
	}()

	if err = deleteStudentTxFunc(tx, userUUID); err != nil {
		return err
	}

	if err = deleteLecturerTxFunc(tx, userUUID); err != nil {
		return err
	}

	if err = deleteUserTxFunc(tx, userUUID); err != nil {
		return err
	}

	return commitTxFunc(tx)
}
