package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, loadErr := util.LoadConfig("../..")
	if loadErr != nil {
		log.Fatal("could not load up the environment config")
	}
	var err error
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to make a connection %v", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
