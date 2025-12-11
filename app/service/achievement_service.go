package service

import (
	"database/sql"
	"encoding/json"
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

// func CreateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
//     form, err := c.MultipartForm()
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{
//             "status":  "error",
//             "message": "invalid_form_data",
//         })
//     }

//    rawStudentID := c.Locals("student_id")
//     studentIDStr, ok := rawStudentID.(string)
//     if !ok || studentIDStr == "" {
//         return c.Status(400).JSON(fiber.Map{
//             "status":  "error",
//             "message": "invalid_or_missing_student_id_in_token",
//         })
//     }

//     studentUUID, err := uuid.Parse(studentIDStr)
//     if err != nil {
//         return c.Status(400).JSON(fiber.Map{
//             "status":  "error",
//             "message": "student_id_not_uuid",
//         })
//     }
//     fmt.Println("STUDENT UUID TOKEN:", studentUUID)

//     getValue := func(key string) string {
//         if arr, ok := form.Value[key]; ok && len(arr) > 0 {
//             return arr[0]
//         }
//         return ""
//     }

//     achievementType := getValue("achievement_type")
//     title := getValue("title")
//     description := getValue("description")
//     status := getValue("status")
//     if status == "" {
//         status = "draft"
//     }

//     pointsStr := getValue("points")
//     points := 0
//     if pointsStr != "" {
//         points, _ = strconv.Atoi(pointsStr)
//     }

//     tags := []string{}
//     if tagsStr := getValue("tags"); tagsStr != "" {
//         json.Unmarshal([]byte(tagsStr), &tags)
//     }

//     var details map[string]interface{}
//     if detailsStr := getValue("details"); detailsStr != "" {
//         json.Unmarshal([]byte(detailsStr), &details)
//     }

//     files := form.File["attachments"]
//     var attachments []model.Attachment

//     for _, file := range files {
//         path := "uploads/" + file.Filename

//         if err := c.SaveFile(file, path); err != nil {
//             return c.Status(500).JSON(fiber.Map{
//                 "status":  "error",
//                 "message": "failed_save_file",
//             })
//         }

//         attachments = append(attachments, model.Attachment{
//             FileName:   file.Filename,
//             FileUrl:    path,
//             FileType:   file.Header.Get("Content-Type"),
//             UploadedAt: time.Now(),
//         })
//     }

//     ach := model.Achievement{
//         StudentID:       studentIDStr,
//         AchievementType: achievementType,
//         Title:           title,
//         Description:     description,
//         Details:         details,
//         Tags:            tags,
//         Points:          points,
//         Attachments:     attachments,
//         Status:          status,
//     }

//     result, err := repository.CreateAchievement(mongoDB, &ach)
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{
//             "status":  "error",
//             "message": "failed_create_achievement_mongo",
//         })
//     }

//     err = repository.CreateAchievementReference(db, studentUUID, result.ID.Hex())
//     if err != nil {
//         return c.Status(500).JSON(fiber.Map{
//             "status":  "error",
//             "message": err.Error(),
//         })
//     }

//     return c.Status(201).JSON(fiber.Map{
//         "status":  "success",
//         "message": "achievement_created_successfully",
//         "data": fiber.Map{
//             "achievement": result,
//         },
//     })
// }

func CreateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {

    req := new(model.CreateAchievement)

    if err := c.BodyParser(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_json",
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
            "message": "failed_create_achievement_mongo",
        })
    }

    studentUUID, _ := uuid.Parse(studentIDStr)
    _ = repository.CreateAchievementReference(db, studentUUID, result.ID.Hex())

    return c.Status(201).JSON(fiber.Map{
        "status": "success",
        "data":   result,
    })
}

func UploadAchievementAttachmentService(c *fiber.Ctx, mongoDB *mongo.Database) error {
    achievementID := c.Params("id")
    if achievementID == "" {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "missing_achievement_id",
        })
    }

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid_form_data",
        })
    }

    files := form.File["attachments"]
    if len(files) == 0 {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "no_files_uploaded",
        })
    }

    if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_create_uploads_folder",
        })
    }

    var attachments []model.Attachment

    for _, file := range files {

        filePath := "uploads/" + file.Filename

        if err := c.SaveFile(file, filePath); err != nil {
            fmt.Println("SAVE FILE ERROR:", err)
            return c.Status(500).JSON(fiber.Map{
                "status":  "error",
                "message": "failed_save_file",
            })
        }

        attachments = append(attachments, model.Attachment{
            FileName:   file.Filename,
            FileUrl:    filePath,
            FileType:   file.Header.Get("Content-Type"),
            UploadedAt: time.Now(),
        })
    }

    // Push ke MongoDB
    err = repository.AddAttachmentsToAchievement(mongoDB, achievementID, attachments)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_attachments",
        })
    }

    return c.JSON(fiber.Map{
        "status":      "success",
        "attachments": attachments,
    })
}


func GetAllAchievementsService(c *fiber.Ctx, mongoDB *mongo.Database, sqlDB *sql.DB) error {
	list, err := repository.GetAllAchievements(mongoDB, sqlDB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_achievements",
			"detail":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_all_achievements",
		"data": fiber.Map{
			"achievements": list,
		},
	})
}

func GetAchievementByIDService(c *fiber.Ctx, mongoDB *mongo.Database, sqlDB *sql.DB) error {
	id := c.Params("id")

	ach, err := repository.GetAchievementByID(mongoDB, sqlDB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_achievement",
			"detail":  err.Error(),
		})
	}

	if ach == nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "achievement_not_found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_achievement",
		"data": fiber.Map{
			"achievement": ach,
		},
	})
}

func UpdateAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id") // mongo achievement id (hex)
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing_id"})
    }

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid_form_data"})
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
        "title":        title,
        "description":  description,
        "details":      details,
        "tags":         tags,
        "attachments":  attachments,
        "status":       status,
        "points":       points,
        "updated_at":   time.Now(),
    })

    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed_update_mongo"})
    }

    err = repository.UpdateAchievementReference(db, id, status)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed_update_reference"})
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement_updated_successfully",
        "data":    updated,
    })
}
func SoftDeleteAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id")

    // Soft delete Mongo
    err := repository.SoftDeleteAchievement(mongoDB, id)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    // Update PostgreSQL reference
    err = repository.UpdateAchievementReference(db, id, "deleted")
    if err != nil {
        msg := err.Error()

        if msg == "reference_not_found" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "achievement_reference_not_found",
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_soft_delete_reference",
            "detail":  msg,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement_deleted_successfully",
    })
}

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
                "message": "achievement_reference_not_found",
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed submit",
            "detail":  msg,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement submitted successfully",
    })
}

func VerifyAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id")

    achievement, err := repository.GetAchievementByID(mongoDB, db, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  "error",
            "message": "achievement_not_found",
        })
    }

    if achievement.Status != "submitted" {
        return c.Status(400).JSON(fiber.Map{
            "status":        "error",
            "message":       "achievement_must_be_submitted_to_verify",
            "current_status": achievement.Status,
        })
    }

    err = repository.UpdateAchievementStatus(mongoDB, id, "verified")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_achievement_mongo",
            "detail":  err.Error(),
        })
    }

    err = repository.UpdateAchievementReference(db, id, "verified")
    if err != nil {
        if err.Error() == "reference_not_found" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "achievement_reference_not_found",
            })
        }
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_reference_postgres",
            "detail":  err.Error(),
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement verified successfully",
    })
}

func RejectAchievementService(c *fiber.Ctx, mongoDB *mongo.Database, db *sql.DB) error {
    id := c.Params("id")

    var body struct {
        RejectionNote string `json:"rejection_note"`
    }
    if err := c.BodyParser(&body); err != nil || strings.TrimSpace(body.RejectionNote) == "" {
        return c.Status(400).JSON(fiber.Map{
            "status":  "error",
            "message": "rejection_note_is_required",
        })
    }

    achievement, err := repository.GetAchievementByID(mongoDB, db, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  "error",
            "message": "achievement_not_found",
        })
    }

    if achievement.Status != "submitted" && achievement.Status != "verified" {
        return c.Status(400).JSON(fiber.Map{
            "status":         "error",
            "message":        "achievement_must_be_submitted_or_verified_to_reject",
            "current_status": achievement.Status,
        })
    }

    err = repository.UpdateAchievementStatus(mongoDB, id, "rejected")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_achievement_status",
            "detail":  err.Error(),
        })
    }

    err = repository.RejectAchievementReference(db, id, body.RejectionNote)
    if err != nil {
        msg := err.Error()
        if msg == "reference_not_found" {
            return c.Status(404).JSON(fiber.Map{
                "status":  "error",
                "message": "achievement_reference_not_found",
            })
        }

        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_update_achievement_reference",
            "detail":  msg,
        })
    }

    return c.Status(200).JSON(fiber.Map{
        "status":  "success",
        "message": "achievement rejected successfully",
    })
}

func GetAchievementsByStudentIDService(c *fiber.Ctx, mongoDB *mongo.Database, sqlDB *sql.DB) error {
	id := c.Params("id")

	achievements, err := repository.GetAchievementsByStudentID(mongoDB, sqlDB, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed_fetch_achievements",
			"detail":  err.Error(),
		})
	}

	if len(achievements) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "no_achievements_found_for_student",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "success_get_achievements",
		"data": fiber.Map{
			"achievements": achievements,
		},
	})
}
func GetAchievementHistoryService(c *fiber.Ctx, mongoDB *mongo.Database, sqlDB *sql.DB) error {
    id := c.Params("id")

    ach, err := repository.GetAchievementByID(mongoDB, sqlDB, id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "status":  "error",
            "message": "failed_fetch_achievement",
            "detail":  err.Error(),
        })
    }
    if ach == nil {
        return c.Status(404).JSON(fiber.Map{
            "status":  "error",
            "message": "achievement_not_found",
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
            updatedBy = ach.VerifiedBy.String()
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
            updatedBy = ach.VerifiedBy.String()
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
        "message": "success_get_achievement_history",
        "data": fiber.Map{
            "history": history,
        },
    })
}

