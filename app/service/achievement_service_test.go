package service

import (
	"database/sql"
	"errors"
	"testing"

	"pelaporan-prestasi/app/model"
	"pelaporan-prestasi/app/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateAchievement(t *testing.T) {
	tests := []struct {
		name      string
		studentID string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:      "Success Create Achievement",
			studentID: uuid.New().String(),
			mockSetup: func() {
				createAchievementFunc = func(db *mongo.Database, a *model.Achievement) (*model.Achievement, error) {
					a.ID = primitive.NewObjectID()
					return a, nil
				}
				createAchievementRefFunc = func(db *sql.DB, id uuid.UUID, achID string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:      "Invalid Student ID",
			studentID: "",
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:      "Mongo Insert Failed",
			studentID: uuid.New().String(),
			mockSetup: func() {
				createAchievementFunc = func(db *mongo.Database, a *model.Achievement) (*model.Achievement, error) {
					return nil, errors.New("mongo error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset
			createAchievementFunc = repository.CreateAchievement
			createAchievementRefFunc = repository.CreateAchievementReference

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req := &model.CreateAchievement{
				AchievementType: "academic",
				Title:           "Juara 1",
				Description:     "Lomba Nasional",
			}

			res, err := CreateAchievement(nil, nil, tt.studentID, req)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && res == nil {
				t.Fatal("expected result")
			}
		})
	}
}
func TestSubmitAchievement(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Success Submit",
			mockSetup: func() {
				updateAchievementStatusFunc = func(*mongo.Database, string, string) error {
					return nil
				}
				updateAchievementRefFunc = func(*sql.DB, string, string) error {
					return nil
				}
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "submitted"}, nil
				}
			},
			wantErr: false,
		},
		{
			name: "Update Mongo Failed",
			mockSetup: func() {
				updateAchievementStatusFunc = func(*mongo.Database, string, string) error {
					return errors.New("mongo error")
				}
			},
			wantErr: true,
		},
		{
			name: "Update Reference Failed",
			mockSetup: func() {
				updateAchievementStatusFunc = func(*mongo.Database, string, string) error {
					return nil
				}
				updateAchievementRefFunc = func(*sql.DB, string, string) error {
					return errors.New("reference error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateAchievementStatusFunc = repository.UpdateAchievementStatus
			updateAchievementRefFunc = repository.UpdateAchievementReference
			getAchievementByIDFunc = repository.GetAchievementByID

			tt.mockSetup()

			_, err := SubmitAchievement(nil, nil, "id")

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestVerifyAchievement(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Success Verify",
			status: "submitted",
			mockSetup: func() {
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "submitted"}, nil
				}
				verifyAchievementFunc = func(*mongo.Database, string, int, string) error {
					return nil
				}
				updateAchievementRefFunc = func(*sql.DB, string, string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:   "Invalid Status",
			status: "draft",
			mockSetup: func() {
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "draft"}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "Verify Mongo Failed",
			mockSetup: func() {
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "submitted"}, nil
				}
				verifyAchievementFunc = func(*mongo.Database, string, int, string) error {
					return errors.New("mongo error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAchievementByIDFunc = repository.GetAchievementByID
			verifyAchievementFunc = repository.VerifyAchievement
			updateAchievementRefFunc = repository.UpdateAchievementReference

			tt.mockSetup()

			_, err := VerifyAchievement(nil, nil, "id", 10, "lecturer")

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestRejectAchievement(t *testing.T) {
	tests := []struct {
		name      string
		status    string
		note      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "Success Reject",
			status: "submitted",
			note:   "Kurang bukti",
			mockSetup: func() {
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "submitted"}, nil
				}
				updateAchievementStatusFunc = func(*mongo.Database, string, string) error {
					return nil
				}
				rejectAchievementRefFunc = func(*sql.DB, string, string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "Empty Note",
			note:    "",
			wantErr: true,
		},
		{
			name:   "Invalid Status",
			status: "draft",
			note:   "Invalid",
			mockSetup: func() {
				getAchievementByIDFunc = func(*mongo.Database, *sql.DB, string) (*model.AchievementResponse, error) {
					return &model.AchievementResponse{Status: "draft"}, nil
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAchievementByIDFunc = repository.GetAchievementByID
			updateAchievementStatusFunc = repository.UpdateAchievementStatus
			rejectAchievementRefFunc = repository.RejectAchievementReference

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			_, err := RejectAchievement(nil, nil, "id", tt.note)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
