package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ridhamu/snippetbox/internal/models"
)

type application struct {
	logger       *slog.Logger
	snippetModel *models.SnippetModels
}

func main() {
	addr := flag.String("addr", ":4000", "HOST:port")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySql Data Source name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// initialized our app struct
	app := application{
		logger: logger,
		snippetModel: &models.SnippetModels{
			DB: db,
		},
	}

	logger.Info(fmt.Sprintf("Server Running on %s", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// using custom middle ware to send 404 to unallowed directory listing
// func neuter(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if strings.HasSuffix(r.URL.Path, "/") {
// 			http.NotFound(w, r)
// 			return
// 		}
//
// 		next.ServeHTTP(w, r)
// 	})
// }

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
