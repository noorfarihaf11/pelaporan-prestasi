package repository

import (
	"database/sql"
	"pelaporan-prestasi/app/model"

	"github.com/google/uuid"
)

func GetAllUser(db *sql.DB) ([]model.User, error) {
	rows, err := db.Query(`SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at 
    FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userList []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userList = append(userList, u)
	}

	return userList, nil
}

func GetUserByID(db *sql.DB, id string) (*model.User, error) {
	row := db.QueryRow(`SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at 
    FROM users WHERE id=$1`, id)

	var u model.User
	err := row.Scan(
			&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	return &u, nil
}

func CreateUserTx(tx *sql.Tx, user *model.User) (*model.User, error) {

    query := `
        INSERT INTO users (username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
    `

    err := tx.QueryRow(
        query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.RoleID,
        user.IsActive,
        user.CreatedAt,
        user.UpdatedAt,
    ).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.RoleID,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }

    return user, nil
}

func UpdateUserTx(tx *sql.Tx, id uuid.UUID, user *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET full_name = $1,
		    username = $2,
		    email = $3,
		    password_hash = $4,
		    role_id = $5,
		    updated_at = NOW()
		WHERE id = $6
		RETURNING id, full_name, username, email, role_id, created_at, updated_at
	`

	err := tx.QueryRow(
		query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.RoleID,
		id,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserTx(tx *sql.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(`DELETE FROM users WHERE id=$1`, userID)
	return err
}

func UpdateUserRole(db *sql.DB, userID uuid.UUID, roleID uuid.UUID) error {
	_, err := db.Exec(`
		UPDATE users
		SET role_id = $1, updated_at = NOW()
		WHERE id = $2
	`, roleID, userID)

	return err
}
