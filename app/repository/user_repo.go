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
		var a model.User
		err := rows.Scan(&a.ID, &a.Username, &a.Email, &a.FullName, &a.RoleID, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userList = append(userList, a)
	}

	return userList, nil
}