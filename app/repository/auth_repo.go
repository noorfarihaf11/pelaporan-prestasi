package repository

import (
	"database/sql"
	"errors"
	"pelaporan-prestasi/app/model"
	_ "time"

	"github.com/google/uuid"
)

func RegisterUser(db *sql.DB, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (full_name, username, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, full_name, username, email, created_at
	`

	err := db.QueryRow(
		query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func LoginUser(db *sql.DB, identifier string) (*model.User, string, error) {
    var user model.User
    var passwordHash string
    var studentID *uuid.UUID // karena student bisa NULL

    query := `
        SELECT 
            u.id,
            u.username,
            u.full_name,
            u.role_id,
            u.password_hash,
            s.id AS student_id
        FROM users u
        LEFT JOIN students s ON s.user_id = u.id
        WHERE u.username = $1
    `

    err := db.QueryRow(query, identifier).Scan(
        &user.ID,
        &user.Username,
        &user.FullName,
        &user.RoleID,
        &passwordHash,
        &studentID,
    )

    if err != nil {
        return nil, "", errors.New("user tidak ditemukan")
    }

    user.StudentID = studentID

    return &user, passwordHash, nil
}

func GetProfile(db *sql.DB, userID string) (*model.User, error) {
	var user model.User

	query := `
		SELECT id, full_name, username, email, created_at
		FROM users
		WHERE id = $1
	`

	err := db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
