package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModels struct {
	DB *sql.DB
}

func (m *UserModels) Insert(name, email, password string) error {
	return nil
}

func (m *UserModels) AuthenticateUser(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModels) Exists(id int) (bool, error) {
	return false, nil
}
