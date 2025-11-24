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
	var req model.CreateUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal parse request",
			"success": false,
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal enkripsi password",
			"success": false,
		})
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "RoleID is not a valid UUID",
			"success": false,
		})
	}

	user := &model.User{
		FullName:     req.FullName,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		 RoleID:       roleID,
		CreatedAt:    time.Now(),
	}

	createdUser, err := repository.CreateUser(db, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user: " + err.Error(),
			"success": false,
		})
	}

	token, err := utils.GenerateToken(*createdUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token JWT",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil membuat user baru",
		"success": true,
		"token":   token,
		"user":    createdUser,
	})
}
