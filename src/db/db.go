package db

import (
	"database/sql"
	"fmt"
	"os"
)

func GetDb() *sql.DB {
	dbHost := os.Getenv("PROPERTY_VIEWER_DB_HOST")
	dbPass := os.Getenv("PROPERTY_VIEWER_DB_PASSWORD")
	if dbPass == "" {
		panic("PROPERTY_VIEWER_DB_PASSWORD must be set.")
	}
    if dbPass == "" {
		panic("PROPERTY_VIEWER_DB_HOST must be set.")
    }

	connString := fmt.Sprintf("property-viewer:%s@tcp(%s:3306)/property_api", dbPass, dbHost)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(180)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
    
	return db
}
