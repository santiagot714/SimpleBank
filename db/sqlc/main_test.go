package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/santiagot714/SimpleBank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	testDB, err = sql.Open("postgres", config.TestDatabaseURL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	testQueries = New(testDB)
	code := m.Run()

	os.Exit(code)
}
