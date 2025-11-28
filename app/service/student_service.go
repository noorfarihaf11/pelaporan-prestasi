package service

import (
	"database/sql"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllStudentService(c *fiber.Ctx, db *sql.DB) error {
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

	userList, err := repository.GetAllStudent(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetching students",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mendapatkan semua data students",
		"data": fiber.Map{
			"students": userList,
		},
	})
}
func GetStudentByIDService(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing student ID",
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

	student, err := repository.GetStudentByID(db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetching student",
		})
	}

	if student == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Student not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mendapatkan data student",
		"data": fiber.Map{
			"student": student,
		},
	})
}