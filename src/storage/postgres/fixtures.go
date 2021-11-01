package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/erp/api/src/config"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var (
	db       *sql.DB
	pgDB     *Postgres
	fixtures *testfixtures.Loader
)

func PrepareFixtures() {
	log.Info("started fixture preparation...")

	err := os.Setenv("POSTGRES_TEST_NAME", "api-db-test")
	if err != nil {
		log.Fatal("failed to set test DB name ", err)
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("failed to load config ", err)
	}

	dbConf := cfg.PostgresTestCfg

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name))
	if err != nil {
		log.Fatal("failed to open test db postgres connect ", err)
	}

	// getting location of fixtures path
	_, filename, _, _ := runtime.Caller(0) // nolint:dogsled
	path := filepath.Dir(filename)

	log.Info("fixture file path ", filepath.Join(path, "data"))

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(filepath.Join(path, "data")),
	)
	if err != nil {
		log.Fatal("failed to create fixtures", err)
	}

	wg := sync.WaitGroup{}

	pgDB, err = New(context.Background(), &wg, &dbConf, &dbConf)
	if err != nil {
		log.Fatal("failed to load test db pg connect ", err)
	}

	log.Info("finished fixture preparation...")
}

// LoadFixtures prepares test db with fixtures.
func LoadFixtures() {
	if err := fixtures.Load(); err != nil {
		log.Fatal("failed to load test db fixtures ", err)
	}
}

func GetDB() *Postgres {
	return pgDB
}
