package service

import (
	"database/sql"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"

	"github.com/gofiber/fiber/v2"
)

// GetAllLecturers godoc
// @Summary      Get all lecturers
// @Description  Mengambil seluruh data dosen (Authorization Bearer)
// @Tags         Lecturer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.BaseResponse{data=map[string][]model.Lecturer}
// @Failure      401 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/lecturers [get]
func GetAllLecturersService(c *fiber.Ctx, db *sql.DB) error {
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

	lecturerList, err := repository.GetAllLecturers(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetching lecturers",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Berhasil mendapatkan semua data lecturers",
		"data": fiber.Map{
			"lecturers": lecturerList,
		},
	})
}
// GetAdviseesByLecturer godoc
// @Summary      Get advisees by lecturer
// @Description  Mengambil daftar mahasiswa bimbingan berdasarkan lecturer ID
// @Tags         Lecturer
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lecturer ID"
// @Success      200  {object}  model.BaseResponse{data=[]model.StudentResponse}
// @Failure      500  {object}  model.BaseResponse
// @Router       /api/v1/lecturers/{id}/advisees [get]
func GetAdviseesByLecturerService(c *fiber.Ctx, db *sql.DB) error {
	lecturerID := c.Params("id")

	advisees, err := repository.GetAdviseesByLecturerID(db, lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_advisees",
			"detail":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_advisees",
		"data":    advisees,
	})
}