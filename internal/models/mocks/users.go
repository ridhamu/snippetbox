package mocks

import (
	"time"

	"github.com/ridhamu/snippetbox/internal/models"
)

type UserModel struct{}

func (u *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}

		return nil
	}

	return models.ErrNoRecord
}

func (u *UserModel) Get(id int) (models.User, error) {
	if id == 1 {
		u := models.User{
			ID:      1,
			Name:    "Alice",
			Email:   "alice@example.com",
			Created: time.Now(),
		}

		return u, nil
	}
	return models.User{}, models.ErrNoRecord
}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDupEmail
	default:
		return nil
	}
}

func (u *UserModel) AuthenticateUser(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (u *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
