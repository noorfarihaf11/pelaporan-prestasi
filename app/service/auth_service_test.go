package service

import (
	"database/sql"
	"testing"

	"pelaporan-prestasi/app/model"

	"github.com/google/uuid"
)

type AuthRepository interface {
	LoginUser(identifier string) (*model.User, string, error)
	RegisterUser(user *model.User) (*model.User, error)
}

func TestLogin_Success(t *testing.T) {
	// simpan fungsi asli
	origLogin := loginUserFunc
	origCheck := checkPasswordFunc
	origToken := generateTokenFunc
	origRefresh := generateRefreshTokenFunc

	defer func() {
		loginUserFunc = origLogin
		checkPasswordFunc = origCheck
		generateTokenFunc = origToken
		generateRefreshTokenFunc = origRefresh
	}()

	// MOCK REPOSITORY
	loginUserFunc = func(db *sql.DB, identifier string) (*model.User, string, error) {
		return &model.User{
			ID:       uuid.New(),
			Username: "noor",
		}, "hashed-password", nil
	}

	// MOCK PASSWORD CHECK
	checkPasswordFunc = func(password, hash string) bool {
		return true
	}

	// MOCK TOKEN
	generateTokenFunc = func(user model.User) (string, error) {
		return "fake-token", nil
	}

	generateRefreshTokenFunc = func(user model.User) (string, error) {
		return "fake-refresh", nil
	}

	user, token, refresh, err := Login(nil, "noor", "password")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "fake-token" || refresh != "fake-refresh" {
		t.Fatalf("unexpected token")
	}

	if user.Username != "noor" {
		t.Fatalf("unexpected user")
	}
}
