package service

import (
	"database/sql"
	"time"

	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"message": "Request body tidak valid",
		})
	}

	userPtr, passwordHash, err := repository.LoginUser(db, req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"message": "Username atau password salah",
		})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"message": "Password salah",
		})
	}

	accessToken, err := utils.GenerateToken(*userPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": "Gagal membuat token JWT",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(*userPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": "Gagal membuat refresh token",
		})
	}

	userResponse := fiber.Map{
		"id":        userPtr.ID,
		"username":  userPtr.Username,
		"full_name":  userPtr.FullName,
		"role":      "Mahasiswa",     
		"permissions": []string{},   
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        accessToken,
			"refreshToken": refreshToken,
			"user":         userResponse,
		},
	})
}


func RegisterService(c *fiber.Ctx, db *sql.DB) error {
	var req model.RegisterRequest
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

	user := &model.User{
		FullName:     req.FullName,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		CreatedAt:	  time.Now(),
	}

	createdUser, err := repository.RegisterUser(db, user)
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
		"message": "Registrasi berhasil",
		"success": true,
		"token":   token,
		"user":    createdUser,
	})
}
