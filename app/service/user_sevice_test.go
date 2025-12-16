package service

import (
	"database/sql"
	"testing"
	_ "time"

	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"
	"pelaporan-prestasi/utils"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		dto         model.UserDTO
		mockSetup   func()
		wantErr     bool
		wantStudent bool
	}{
		{
			name: "Valid Create User with Student",
			dto: model.UserDTO{
				FullName: "Noor",
				Username: "noor",
				Email:    "noor@gmail.com",
				Password: "password",
				RoleID:   uuid.New().String(),
				StudentProfile: &model.StudentProfileDTO{
					StudentID:    "123",
					ProgramStudy: "TI",
					AcademicYear: "2023",
				},
			},
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) {
					return &sql.Tx{}, nil
				}
				hashPasswordFunc = func(p string) (string, error) {
					return "hashed", nil
				}
				createUserTxFunc = func(tx *sql.Tx, u *model.User) (*model.User, error) {
					u.ID = uuid.New()
					return u, nil
				}
				createStudentTxFunc = func(tx *sql.Tx, s *model.Student) error {
					return nil
				}
				commitTxFunc = func(tx *sql.Tx) error { return nil }
				rollbackTxFunc = func(tx *sql.Tx) error { return nil }
			},
			wantErr:     false,
			wantStudent: true,
		},
		{
			name: "Invalid Role ID",
			dto: model.UserDTO{
				FullName: "Noor",
				Password: "password",
				RoleID:   "invalid-uuid",
			},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name: "Hash Password Failed",
			dto: model.UserDTO{
				FullName: "Noor",
				Password: "password",
				RoleID:   uuid.New().String(),
			},
			mockSetup: func() {
				hashPasswordFunc = func(p string) (string, error) {
					return "", sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name: "Create User Failed",
			dto: model.UserDTO{
				FullName: "Noor",
				Password: "password",
				RoleID:   uuid.New().String(),
			},
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) {
					return &sql.Tx{}, nil
				}
				hashPasswordFunc = func(p string) (string, error) {
					return "hashed", nil
				}
				createUserTxFunc = func(tx *sql.Tx, u *model.User) (*model.User, error) {
					return nil, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset mock ke default
			beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
			hashPasswordFunc = utils.HashPassword
			createUserTxFunc = repository.CreateUserTx
			createStudentTxFunc = repository.CreateStudentTx
			commitTxFunc = func(tx *sql.Tx) error { return nil }
			rollbackTxFunc = func(tx *sql.Tx) error { return nil }

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			user, student, _, err := CreateUser(nil, tt.dto)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user == nil {
				t.Fatal("expected user, got nil")
			}

			if tt.wantStudent && student == nil {
				t.Fatal("expected student profile")
			}
		})
	}
}
func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		dto       model.UserDTO
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Valid Update",
			userID: uuid.New().String(),
			dto: model.UserDTO{
				FullName: "Updated Name",
			},
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
				commitTxFunc = func(tx *sql.Tx) error { return nil }
				rollbackTxFunc = func(tx *sql.Tx) error { return nil }

				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return &model.User{ID: uuid.New()}, nil
				}

				updateUserTxFunc = func(tx *sql.Tx, id uuid.UUID, u *model.User) (*model.User, error) {
					return u, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "Invalid User ID",
			userID:    "invalid-uuid",
			dto:       model.UserDTO{FullName: "Updated"},
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:   "User Not Found",
			userID: uuid.New().String(),
			dto:    model.UserDTO{FullName: "Updated"},
			mockSetup: func() {
				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return nil, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name:   "Update User Failed",
			userID: uuid.New().String(),
			dto:    model.UserDTO{FullName: "Updated"},
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return &model.User{ID: uuid.New()}, nil
				}
				updateUserTxFunc = func(tx *sql.Tx, id uuid.UUID, u *model.User) (*model.User, error) {
					return nil, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset mock ke default
			beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
			commitTxFunc = func(tx *sql.Tx) error { return nil }
			rollbackTxFunc = func(tx *sql.Tx) error { return nil }
			getUserByIDFunc = repository.GetUserByID
			updateUserTxFunc = repository.UpdateUserTx

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			user, _, _, err := UpdateUser(nil, tt.userID, tt.dto)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && user == nil {
				t.Fatal("expected updated user")
			}
		})
	}
}
func TestDeleteUser(t *testing.T) {
	validID := uuid.New().String()

	tests := []struct {
		name      string
		userID    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Valid Delete",
			userID: validID,
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
				commitTxFunc = func(tx *sql.Tx) error { return nil }
				rollbackTxFunc = func(tx *sql.Tx) error { return nil }

				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return &model.User{ID: uuid.MustParse(id)}, nil
				}

				deleteStudentTxFunc = func(tx *sql.Tx, id uuid.UUID) error { return nil }
				deleteLecturerTxFunc = func(tx *sql.Tx, id uuid.UUID) error { return nil }
				deleteUserTxFunc = func(tx *sql.Tx, id uuid.UUID) error { return nil }
			},
			wantErr: false,
		},
		{
			name:      "Invalid User ID",
			userID:    "invalid-uuid",
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:   "User Not Found",
			userID: validID,
			mockSetup: func() {
				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return nil, sql.ErrNoRows
				}
			},
			wantErr: true,
		},
		{
			name:   "Delete User Failed",
			userID: validID,
			mockSetup: func() {
				beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return &sql.Tx{}, nil }

				getUserByIDFunc = func(db *sql.DB, id string) (*model.User, error) {
					return &model.User{ID: uuid.MustParse(id)}, nil
				}

				deleteStudentTxFunc = func(tx *sql.Tx, id uuid.UUID) error {
					return nil
				}
				deleteLecturerTxFunc = func(tx *sql.Tx, id uuid.UUID) error {
					return nil
				}
				deleteUserTxFunc = func(tx *sql.Tx, id uuid.UUID) error {
					return sql.ErrNoRows
				}

				commitTxFunc = func(tx *sql.Tx) error { return nil }
				rollbackTxFunc = func(tx *sql.Tx) error { return nil }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset mock
			beginTxFunc = func(db *sql.DB) (*sql.Tx, error) { return nil, nil }
			commitTxFunc = func(tx *sql.Tx) error { return nil }
			rollbackTxFunc = func(tx *sql.Tx) error { return nil }
			getUserByIDFunc = repository.GetUserByID
			deleteStudentTxFunc = repository.DeleteStudentTx
			deleteLecturerTxFunc = repository.DeleteLecturerTx
			deleteUserTxFunc = repository.DeleteUserTx

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := DeleteUser(nil, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
