package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	logger *slog.Logger
}

func main() {

	// command line arguments passed it at app runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "file:data.db", "SQLite3 data source name")
	flag.Parse()
	// create a new logger by passing a logger handler
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// create instance of the application struct
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
