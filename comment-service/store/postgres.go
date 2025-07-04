package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func Init() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		if host == "" || port == "" || user == "" || name == "" {
			log.Fatal("one of DB_HOST, DB_PORT, DB_USER, DB_NAME is not set")
		}
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, pass, name,
		)
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	DB = db
	migrate()
}

func migrate() {
	q := `CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        news_id INT NOT NULL,
        parent_id INT,
        text TEXT NOT NULL,
        created TIMESTAMP NOT NULL DEFAULT now()
    );`
	if _, err := DB.Exec(q); err != nil {
		log.Fatalf("migrate comments: %v", err)
	}
}
