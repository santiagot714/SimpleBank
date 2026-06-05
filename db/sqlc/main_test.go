package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
var testStore *Store

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	testQueries = New(testDB)
	testStore = NewStore(testDB)
	code := m.Run()
	// defer func() {
	// 	err = testDB.Close()
	// 	if err != nil {
	// 		log.Fatal("Cannot close database: ", err)
	// 	}
	// }()
	os.Exit(code)
}
