package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ridhamu/snippetbox/internal/models"
)

type application struct {
	logger         *slog.Logger
	snippetModel   *models.SnippetModels
	userModel      *models.UserModels
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HOST:port")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySql Data Source name")

	flag.Parse()

	// INIT ALL REQUIRED DEPENDENCIES
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	cacheTemplate, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		_ = db.Close()
	}()

	// initialize form decoder
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// app init
	app := application{
		logger: logger,
		snippetModel: &models.SnippetModels{
			DB: db,
		},
		userModel: &models.UserModels{
			DB: db,
		},
		templateCache:  cacheTemplate,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// tls config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// init server configuration
	srv := http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// END INIT

	// STARTING SERVER
	logger.Info(fmt.Sprintf("Server Running on %s", *addr))

	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("./tls/localhost+1.pem", "./tls/localhost+1-key.pem")

	logger.Error(err.Error())
	os.Exit(1)
	// END STARTING SERVER
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
