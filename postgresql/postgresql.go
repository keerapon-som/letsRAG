package postgresql

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	// Optionally, you can ping the database to ensure the connection is established
	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

// DB returns the database connection
func DB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection is not initialized")
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
