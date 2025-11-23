package repository

import (
	"database/sql"
	"errors"
	"pelaporan-prestasi/app/model"
	_ "time"
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

	query := `
		SELECT id, username, email, password_hash, full_name, created_at
		FROM users
		WHERE username = $1 OR email = $1
	`

	err := db.QueryRow(query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.FullName,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, "", errors.New("user tidak ditemukan")
	}

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
