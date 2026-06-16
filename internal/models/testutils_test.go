package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	dbPool, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// setup the tables
	setupScript, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		dbPool.Close()
		t.Fatal(err)
	}

	_, err = dbPool.Exec(string(setupScript))
	if err != nil {
		dbPool.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer dbPool.Close()
		teardownScript, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = dbPool.Exec(string(teardownScript))
		if err != nil {
			t.Fatal(err)
		}
	})

	return dbPool
}
