package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "fmt"
	"os"
	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateAchievement godoc
// @Summary      Create achievement
// @Description  Mahasiswa membuat achievement (status awal: draft)
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      model.CreateAchievement  true  "Create Achievement Payload"
// @Success      201 {object} model.BaseResponse{data=model.AchievementSwagger}
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements [post]
func CreateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {

	req := new(model.CreateAchievement)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON payload",
		})
	}

	rawStudentID := c.Locals("student_id")
	studentIDStr, ok := rawStudentID.(string)
	if !ok || studentIDStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid_or_missing_student_id",
		})
	}

	ach := model.Achievement{
		StudentID:       studentIDStr,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Status:          "draft",
		Attachments:     []model.Attachment{},
	}

	result, err := repository.CreateAchievement(mongoDB, &ach)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed create achievement mongo",
		})
	}

	studentUUID, _ := uuid.Parse(studentIDStr)
	_ = repository.CreateAchievementReference(db, studentUUID, result.ID.Hex())

	return c.Status(201).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

// UploadAchievementAttachment godoc
// @Summary      Upload achievement attachments
// @Description  Upload file pendukung achievement
// @Tags         Achievement
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id           path      string true  "Achievement ID"
// @Param        attachments  formData  file   true  "Attachment files"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/attachments [post]
func UploadAchievementAttachmentService(c *fiber.Ctx, mongoDB *mongo.Database) error {
	achievementID := c.Params("id")
	if achievementID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing achievement id",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid form data",
		})
	}

	files := form.File["attachments"]
	if len(files) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "No files upload",
		})
	}

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed creates upload folder",
		})
	}

	var attachments []model.Attachment

	for _, file := range files {

		filePath := "uploads/" + file.Filename

		if err := c.SaveFile(file, filePath); err != nil {
			fmt.Println("SAVE FILE ERROR:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed save file",
			})
		}

		attachments = append(attachments, model.Attachment{
			FileName:   file.Filename,
			FileUrl:    filePath,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		})
	}

	err = repository.AddAttachmentsToAchievement(mongoDB, achievementID, attachments)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update attachments",
		})
	}

	return c.JSON(fiber.Map{
		"status":      "success",
		"attachments": attachments,
	})
}

// // GetAllAchievements godoc
// // @Summary      Get all achievements
// // @Description  Mengambil semua achievement (admin)
// // @Tags         Achievement
// // @Produce      json
// // @Security     BearerAuth
// // @Success      200 {object} model.BaseResponse
// // @Failure      500 {object} model.BaseResponse
// // @Router       /api/v1/achievements [get]
// func GetAllAchievementsService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
// 	list, err := repository.GetAllAchievements(mongoDB, db)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed fetch attachments",
// 			"detail":  err.Error(),
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "Success get all achievements",
// 		"data": fiber.Map{
// 			"achievements": list,
// 		},
// 	})
// }

// GetAchievements godoc
// @Summary      Get achievements by role
// @Description  Mengambil achievement berdasarkan role (Admin / Mahasiswa / Dosen Wali)
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        page          query     int     false "Page number"
// @Param        limit         query     int     false "Limit data"
// @Param        sort          query     string  false "Sort field"
// @Param        order         query     string  false "ASC / DESC"
// @Param        status        query     string  false "Achievement status"
// @Param        student_name  query     string  false "Student name"
// @Success      200 {object} model.BaseResponse
// @Failure      401 {object} model.BaseResponse
// @Failure      403 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements [get]
func GetAchievementsService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    rawUserID := c.Locals("user_id")
    rawRoleID := c.Locals("role_id")

    userID, ok1 := rawUserID.(string)
    roleID, ok2 := rawRoleID.(string)

    if !ok1 || !ok2 || userID == "" || roleID == "" {
        return c.Status(401).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid token or missing role id",
        })
    }

    roleName, err := repository.GetRoleNameByID(db, roleID)
    if err != nil {
        return c.Status(403).JSON(fiber.Map{
            "status":  "error",
            "message": "Cannot determine role name",
        })
    }

    var (
        achievements []model.AchievementResponse
        total        int
    )

    switch roleName {

    case "Admin":
        limit, _ := strconv.Atoi(c.Query("limit", "10"))
        page, _ := strconv.Atoi(c.Query("page", "1"))
        offset := (page - 1) * limit

        sort := c.Query("sort", "created_at")
        order := strings.ToUpper(c.Query("order", "DESC"))
        status := c.Query("status", "")
        studentName := c.Query("student_name", "")

        achievements, total, err = repository.GetAdminAchievementsPaginated(
            db, mongoDB, limit, offset, sort, order, status, studentName,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status": "error",
                "message": "Failed fetch attachment",
                "detail": err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Success get achievement",
            "data": fiber.Map{
                "total":        total,
                "page":         page,
                "limit":        limit,
                "achievements": achievements,
            },
        })

    case "Mahasiswa":
        studentID, err := repository.GetStudentIDByUserID(db, userID)
        if err != nil || studentID == "" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "Student not found",
            })
        }

        achievements, err = repository.GetAchievementsForStudent(
            mongoDB, db, studentID,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status": "error",
                "message": "Failed fetch achievement",
                "detail": err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Success get achievement",
            "data": fiber.Map{
                "achievements": achievements,
            },
        })

    case "Dosen Wali":
        lecturerID, err := repository.GetLecturerIDByUserID(db, userID)
        if err != nil || lecturerID == "" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "Lecturer not found",
            })
        }

        limit, _ := strconv.Atoi(c.Query("limit", "10"))
        page, _ := strconv.Atoi(c.Query("page", "1"))
        offset := (page - 1) * limit

        achievements, total, err = repository.GetLecturerAchievementsPaginated(
            db, mongoDB, lecturerID, limit, offset,
        )
        if err != nil {
            return c.Status(500).JSON(fiber.Map{
                "status": "error",
                "message": "Failed fetch achievement",
                "detail": err.Error(),
            })
        }

        return c.JSON(fiber.Map{
            "status":  "success",
            "message": "Success get achievement",
            "data": fiber.Map{
                "total":        total,
                "page":         page,
                "limit":        limit,
                "achievements": achievements,
            },
        })

    default:
        return c.Status(403).JSON(fiber.Map{
            "status":  "error",
            "message": "Forbidden invalid role",
        })
    }
}


// GetAchievementByID godoc
// @Summary      Get achievement detail
// @Description  Mengambil detail achievement berdasarkan ID
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Achievement ID"
// @Success      200 {object} model.BaseResponse{data=map[string]model.AchievementSwagger}
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id} [get]
func GetAchievementByIDService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	ach, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetch achievements",
			"detail":  err.Error(),
		})
	}

	if ach == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Achievement not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success get achievement",
		"data": fiber.Map{
			"achievement": ach,
		},
	})
}

// UpdateAchievement godoc
// @Summary      Update achievement
// @Description  Update achievement (multipart form, termasuk attachment)
// @Tags         Achievement
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id           path      string true  "Achievement ID"
// @Param        title        formData  string false "Title"
// @Param        description  formData  string false "Description"
// @Param        status       formData  string false "Status"
// @Param        details      formData  string false "JSON details"
// @Param        tags         formData  string false "JSON tags"
// @Param        points       formData  int    false "Points"
// @Param        attachments  formData  file   false "Attachments"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id} [put]
func UpdateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id") 
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Missing ID"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid form data"})
	}

	title := form.Value["title"][0]
	description := form.Value["description"][0]
	status := "draft"
	if len(form.Value["status"]) > 0 {
		status = form.Value["status"][0]
	}

	var details map[string]interface{}
	if len(form.Value["details"]) > 0 {
		json.Unmarshal([]byte(form.Value["details"][0]), &details)
	}

	tags := []string{}
	if len(form.Value["tags"]) > 0 {
		json.Unmarshal([]byte(form.Value["tags"][0]), &tags)
	}

	files := form.File["attachments"]
	var attachments []model.Attachment

	for _, file := range files {
		path := "uploads/" + file.Filename
		c.SaveFile(file, path)

		attachments = append(attachments, model.Attachment{
			FileName:   file.Filename,
			FileUrl:    path,
			FileType:   file.Header.Get("Content-Type"),
			UploadedAt: time.Now(),
		})
	}

	points := 0
	if len(form.Value["points"]) > 0 {
		points, _ = strconv.Atoi(form.Value["points"][0])
	}

	updated, err := repository.UpdateAchievement(mongoDB, id, bson.M{
		"title":       title,
		"description": description,
		"details":     details,
		"tags":        tags,
		"attachments": attachments,
		"status":      status,
		"points":      points,
		"updated_at":  time.Now(),
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed update mongo"})
	}

	err = repository.UpdateAchievementReference(db, id, status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed update reference"})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement updated successfully",
		"data":    updated,
	})
}

// SoftDeleteAchievement godoc
// @Summary      Soft delete achievement
// @Description  Menghapus achievement (soft delete)
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id} [delete]
func SoftDeleteAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	err := repository.SoftDeleteAchievement(mongoDB, id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	err = repository.UpdateAchievementReference(db, id, "deleted")
	if err != nil {
		msg := err.Error()

		if msg == "reference_not_found" {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Achievement reference not found",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed soft delete reference",
			"detail":  msg,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement deleted successfully",
	})
}

// SubmitAchievement godoc
// @Summary      Submit achievement
// @Description  Mengirim achievement (draft → submitted)
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/submit [post]
func SubmitAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	err := repository.UpdateAchievementStatus(mongoDB, id, "submitted")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	err = repository.UpdateAchievementReference(db, id, "submitted")
	if err != nil {
		msg := err.Error()

		if msg == "reference_not_found" {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Achievement reference not found",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed submit",
			"detail":  msg,
		})
	}

	updated, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed get updated achievement",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "achievement submitted successfully",
		"data": fiber.Map{
			"achievement": updated,
		},
	})
}

// VerifyAchievement godoc
// @Summary      Verify achievement
// @Description  Verifikasi achievement oleh dosen wali
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                  true "Achievement ID"
// @Param        body  body      model.VerifyAchievement true "Verification payload"
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/verify [post]
func VerifyAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	achievement, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil || achievement == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Achievement not found",
		})
	}

	if achievement.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"status":         "error",
			"message":        "Achievement must be verify",
			"current_status": achievement.Status,
		})
	}

	var payload model.VerifyAchievement
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	rawUserID := c.Locals("user_id")
	lecturerID, _ := rawUserID.(string)

	err = repository.VerifyAchievement(mongoDB, id, payload.Points, lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update achievement mongo",
			"detail":  err.Error(),
		})
	}

	err = repository.UpdateAchievementReference(db, id, "verified")
	if err != nil {
		if err.Error() == "Reference not found" {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Achievement reference not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed updated reference postgres",
			"detail":  err.Error(),
		})
	}

	updated, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed get update achievement",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "achievement verified successfully",
		"data": fiber.Map{
			"achievement":    updated,
			"updated_status": updated.Status, 
		},
	})
}

// RejectAchievement godoc
// @Summary      Reject achievement
// @Description  Menolak achievement dengan catatan
// @Tags         Achievement
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string true "Achievement ID"
// @Param        body  body      object true "Rejection payload" example({"rejection_note":"Data tidak valid"})
// @Success      200 {object} model.BaseResponse
// @Failure      400 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/reject [post]
func RejectAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	var body struct {
		RejectionNote string `json:"rejection_note"`
	}
	if err := c.BodyParser(&body); err != nil || strings.TrimSpace(body.RejectionNote) == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Rejection note required",
		})
	}

	achievement, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Achievement not found",
		})
	}

	if achievement.Status != "submitted" && achievement.Status != "verified" {
		return c.Status(400).JSON(fiber.Map{
			"status":         "error",
			"message":        "Achievement must be submitted or verified to reject.",
			"current_status": achievement.Status,
		})
	}

	err = repository.UpdateAchievementStatus(mongoDB, id, "rejected")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update achievement status",
			"detail":  err.Error(),
		})
	}

	err = repository.RejectAchievementReference(db, id, body.RejectionNote)
	if err != nil {
		msg := err.Error()
		if msg == "reference_not_found" {
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Achievement reference not found",
			})
		}

		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed update achievement references",
			"detail":  msg,
		})
	}

	rejected, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed get rejected achievements",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "achievement rejected successfully",
		"data": fiber.Map{
			"achievement":     rejected,
			"rejected_status": rejected.Status,
		},
	})
}

// GetAchievementsByStudentID godoc
// @Summary      Get achievements by student
// @Description  Mengambil achievement berdasarkan student ID
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string true "Student ID"
// @Success      200 {object} model.BaseResponse
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/student [get]
func GetAchievementsByStudentIDService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	achievements, err := repository.GetAchievementsByStudentID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetch achievements",
			"detail":  err.Error(),
		})
	}

	if len(achievements) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No achievements found for student",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success get achievement",
		"data": fiber.Map{
			"achievements": achievements,
		},
	})
}

// GetAchievementHistory godoc
// @Summary      Get achievement history
// @Description  Riwayat perubahan status achievement
// @Tags         Achievement
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string true "Achievement ID"
// @Success      200 {object} model.BaseResponse{data=map[string][]model.AchievementHistory}
// @Failure      404 {object} model.BaseResponse
// @Failure      500 {object} model.BaseResponse
// @Router       /api/v1/achievements/{id}/history [get]
func GetAchievementHistoryService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
	id := c.Params("id")

	ach, err := repository.GetAchievementByID(mongoDB, db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed fetch achievement",
			"detail":  err.Error(),
		})
	}
	if ach == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Achievement not found",
		})
	}

	var history []model.AchievementHistory

	// Draft
	history = append(history, model.AchievementHistory{
		Status:    "draft",
		UpdatedAt: ach.CreatedAt,
		UpdatedBy: ach.StudentID,
	})

	// Submitted
	if ach.Status != "draft" {
		history = append(history, model.AchievementHistory{
			Status:    "submitted",
			UpdatedAt: ach.UpdatedAt,
			UpdatedBy: ach.StudentID,
		})
	}

	// Verified
	if ach.VerifiedAt != nil {
		updatedBy := ""
		if ach.VerifiedBy != nil {
			updatedBy = *ach.VerifiedBy

		}
		history = append(history, model.AchievementHistory{
			Status:        "verified",
			UpdatedAt:     *ach.VerifiedAt,
			UpdatedBy:     updatedBy,
			RejectionNote: ach.RejectionNote,
		})
	}

	// Rejected
	if ach.Status == "rejected" && ach.RejectionNote != nil {
		updatedBy := ""
		if ach.VerifiedBy != nil {
			updatedBy = *ach.VerifiedBy
		}
		updatedAt := time.Now()
		if ach.VerifiedAt != nil {
			updatedAt = *ach.VerifiedAt
		}
		history = append(history, model.AchievementHistory{
			Status:        "rejected",
			UpdatedAt:     updatedAt,
			UpdatedBy:     updatedBy,
			RejectionNote: ach.RejectionNote,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success get achievement history",
		"data": fiber.Map{
			"history": history,
		},
	})
}

//unit testing
var (
	createAchievementFunc    = repository.CreateAchievement
	createAchievementRefFunc = repository.CreateAchievementReference
)
var (
	updateAchievementStatusFunc    = repository.UpdateAchievementStatus
	updateAchievementRefFunc       = repository.UpdateAchievementReference
	getAchievementByIDFunc         = repository.GetAchievementByID
	verifyAchievementFunc          = repository.VerifyAchievement
	rejectAchievementRefFunc       = repository.RejectAchievementReference
)

// test logic, karna fiber tidak bisa dilakukan unit testing
func CreateAchievement(
	mongoDB *mongo.Database,
	db *sql.DB,
	studentID string,
	req *model.CreateAchievement,
) (*model.Achievement, error) {

	if studentID == "" {
		return nil, errors.New("invalid student id")
	}

	ach := &model.Achievement{
		StudentID:       studentID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         req.Details,
		Tags:            req.Tags,
		Status:          "draft",
		Attachments:     []model.Attachment{},
	}

	result, err := createAchievementFunc(mongoDB, ach)
	if err != nil {
		return nil, err
	}

	studentUUID, _ := uuid.Parse(studentID)
	_ = createAchievementRefFunc(db, studentUUID, result.ID.Hex())

	return result, nil
}
func SubmitAchievement(
	mongoDB *mongo.Database,
	db *sql.DB,
	id string,
) (*model.AchievementResponse, error) {

	if err := updateAchievementStatusFunc(mongoDB, id, "submitted"); err != nil {
		return nil, err
	}

	if err := updateAchievementRefFunc(db, id, "submitted"); err != nil {
		return nil, err
	}

	return getAchievementByIDFunc(mongoDB, db, id)
}
func VerifyAchievement(
	mongoDB *mongo.Database,
	db *sql.DB,
	id string,
	points int,
	lecturerID string,
) (*model.AchievementResponse, error) {

	ach, err := getAchievementByIDFunc(mongoDB, db, id)
	if err != nil || ach == nil {
		return nil, errors.New("not found")
	}

	if ach.Status != "submitted" {
		return nil, errors.New("invalid status")
	}

	if err := verifyAchievementFunc(mongoDB, id, points, lecturerID); err != nil {
		return nil, err
	}

	if err := updateAchievementRefFunc(db, id, "verified"); err != nil {
		return nil, err
	}

	return getAchievementByIDFunc(mongoDB, db, id)
}
func RejectAchievement(
	mongoDB *mongo.Database,
	db *sql.DB,
	id string,
	note string,
) (*model.AchievementResponse, error) {

	if strings.TrimSpace(note) == "" {
		return nil, errors.New("rejection note required")
	}

	ach, err := getAchievementByIDFunc(mongoDB, db, id)
	if err != nil {
		return nil, err
	}

	if ach.Status != "submitted" && ach.Status != "verified" {
		return nil, errors.New("invalid status")
	}

	if err := updateAchievementStatusFunc(mongoDB, id, "rejected"); err != nil {
		return nil, err
	}

	if err := rejectAchievementRefFunc(db, id, note); err != nil {
		return nil, err
	}

	return getAchievementByIDFunc(mongoDB, db, id)
}
