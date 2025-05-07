package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"dream.website/internal/model"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       model.SnippetModel
	Users          model.UserModelInterface
	Templatecache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}

func main() {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		log.Fatal("APP_ADDR is not set")
	}
	log.Printf("Starting server on %s", addr)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := OpenDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	Templatecache, err := NewTemplatecache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()

	SessionManager := scs.New()
	SessionManager.Store = mysqlstore.New(db)
	SessionManager.Lifetime = 12 * time.Hour
	SessionManager.Cookie.Secure = true

	app := &Application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		Users:          &model.UserModel{DB: db},
		snippets:       &model.SnipppetModel{DB: db},
		Templatecache:  Templatecache,
		FormDecoder:    formDecoder,
		SessionManager: SessionManager,
	}

	// Load CA certificate
	rootCertPool := x509.NewCertPool()
	caPem, err := os.ReadFile(os.Getenv("CA_PEM_PATH"))
	if err != nil {
		errorLog.Fatal("Unable to read CA file:", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(caPem); !ok {
		errorLog.Fatal("Failed to append CA certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:          rootCertPool,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServeTLS(os.Getenv("CERT_PEM_PATH"), os.Getenv("KEY_PEM_PATH"))
	errorLog.Fatal(err)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
