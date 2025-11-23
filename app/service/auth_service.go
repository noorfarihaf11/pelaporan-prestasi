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
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request",
		})
	}

	userPtr, passwordHash, err := repository.LoginUser(db, req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Credentials",
		})
	}

	token, err := utils.GenerateToken(*userPtr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(*userPtr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token":        token,
			"refreshToken": refreshToken,
			"user": fiber.Map{
				"id":          userPtr.ID,
				"full_name":   userPtr.FullName,
				"username":    userPtr.Username,
				"role":        "Mahasiswa",
				"permissions": []string{},
			},
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
		CreatedAt:    time.Now(),
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

func RefreshTokenService(c *fiber.Ctx, db *sql.DB) error {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request",
		})
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Refresh Token",
		})
	}

	username := claims["username"].(string)

	userPtr, _, err := repository.LoginUser(db, username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User Not Found",
		})
	}

	newAccessToken, err := utils.GenerateToken(*userPtr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"token": newAccessToken,
		},
	})
}
func LogoutService(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Logout berhasil",
	})
}
