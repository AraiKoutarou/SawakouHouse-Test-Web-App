// shared/db: アプリ全体で共有するDB接続を管理します。
// テーブルの作成など、モジュール固有の処理はここには書きません。
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connect: 環境変数からDB接続情報を読み取り、接続済みの *sql.DB を返します。
func Connect() *sql.DB {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "sawakou_board")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("DBへのping失敗: %v", err)
	}

	log.Println("DB接続成功")
	return db
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
