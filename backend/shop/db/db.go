package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "sawakou_board")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("DBへのping失敗: %v", err)
	}

	log.Println("DB接続成功")
	createTable()
}

func createTable() {
	query_shop := `
	CREATE TABLE IF NOT EXISTS posts_shop (
		id         	SERIAL PRIMARY KEY,
		name      	VARCHAR(255) NOT NULL,
		price     	INT NOT NULL,
		description VARCHAR(100) NOT NULL,
		created_at 	TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	if _, err := DB.Exec(query_shop); err != nil {
		log.Fatalf("テーブル作成失敗: %v", err)
	}
	log.Println("posts_shopテーブル確認完了")
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
