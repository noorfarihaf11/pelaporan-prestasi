package middleware

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"pelaporan-prestasi/utils"
)

func RBAC(permission string, db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "missing_token",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid_token_format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid_token",
			})
		}

		roleID := claims.RoleID

		if roleID == uuid.Nil {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "missing_role_claim",
			})
		}

		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT 1 
				FROM role_permissions rp
				JOIN permissions p ON p.id = rp.permission_id
				WHERE rp.role_id = $1
				AND p.name = $2
			)
		`, roleID, permission).Scan(&exists)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "db_error_permission_check",
			})
		}

		if !exists {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "forbidden",
			})
		}

		return c.Next()
	}
}
