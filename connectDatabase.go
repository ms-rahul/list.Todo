package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func importFromEnv(key string) string {
	err := godotenv.Load("env_var.env")

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func createTable(db *sql.DB) {
	q := "CREATE TABLE IF NOT EXISTS todo_list(id int primary key auto_increment, todo_desc varchar(300), status bool)"
	_, err := db.Exec(q)
	if err != nil {
		log.Fatal(err)
	}

}

func createAndUseDB() *sql.DB {
	host := importFromEnv("DB_HOSTNAME")
	username := importFromEnv("DB_USERNAME")
	password := importFromEnv("DB_PASSWORD")
	database := importFromEnv("DB_NAME")

	db_creds := username + ":" + password + "@" + host + "/"
	db, err := sql.Open("mysql", db_creds)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + database)
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("mysql", db_creds+database)
	if err != nil {
		log.Fatal(err)
	}

	// create table to hold to-do
	createTable(db)

	return db
}
