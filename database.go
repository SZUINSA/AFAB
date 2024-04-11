package main

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/sqlite"
	"log"
	"os"
)

var db *sql.DB

func init_db() error {
	var err error
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0755)
		if err != nil {
			fmt.Println("Fail to create data folder", err)
			return err
		}
	}
	db, err = sql.Open("sqlite", "data/sqlite.db")
	if err != nil {
		log.Println("Failed to open database: ", err)
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return err
	}
	log.Println("Connected to database ")
	var stmt *sql.Stmt
	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS 'game' (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
			GroupId TEXT NOT NULL,
			FolderToken TEXT NOT NULL,
			FileToken TEXT NOT NULL
        )
	`)
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Println("Table 'game' create failed", err)
		return err
	}
	log.Println("Table 'game' is created successfully.")
	return err
}
