package repository

import (
	"database/sql"
	"pelaporan-prestasi/app/model"
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


func UpdateUser(db *sql.DB, id string, user *model.User) (*model.User, error) {
	query := `
		UPDATE (full_name, username, email, password_hash, role_id)
		SET ($1, $2, $3, $4, $5)
		RETURNING id, full_name, username, email, role_id, created_at
	`

	err := db.QueryRow(
		query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.RoleID,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.Email,
		&user.RoleID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
