package models

import (
	"testing"

	"github.com/ridhamu/snippetbox/internal/assert"
)

func TestUserExist(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dbPool := newTestDB(t)

			UserModelTest := UserModel{
				DB: dbPool,
			}
			isExist, err := UserModelTest.Exists(test.userID)
			assert.Equal(t, isExist, test.want)
			assert.NilError(t, err)
		})
	}
}
