package main

import (
	"database/sql"
	"flag"

	//new

	"github.com/po133na/go-mid/pkg/jsonlog"
	"github.com/po133na/go-mid/pkg/shelter/model"

	//new
	"sync"

	_ "github.com/lib/pq"

	//new
	"os"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:ddoossyymm1211@localhost:5433/Animal_Shelter?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()
	//new
	logger := jsonlog.NewLogger(os.Stdout, jsonlog.LevelInfo)

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
		//new
		logger: logger,
	}
	if err := app.serve(); err != nil {
		logger.PrintFatal(err, nil)
	}

	app.routes()
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config // struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
