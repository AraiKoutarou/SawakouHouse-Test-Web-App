// internal/infrastructure/persistence/post/repository.go: DB実装。
package post

import (
	"database/sql"
	"log"

	"github.com/arakou0812/backend/internal/domain/post"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id         SERIAL PRIMARY KEY,
		title      VARCHAR(255) NOT NULL,
		content    TEXT NOT NULL,
		author     VARCHAR(100) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	if _, err := r.db.Exec(query); err != nil {
		log.Fatalf("postsテーブル作成失敗: %v", err)
	}
	log.Println("postsテーブル確認完了")
}

func (r *PostgresRepository) GetAll() ([]post.Post, error) {
	rows, err := r.db.Query(`SELECT id, title, content, author, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []post.Post{}
	}
	return posts, nil
}

func (r *PostgresRepository) GetByID(id int) (post.Post, error) {
	var p post.Post
	err := r.db.QueryRow(
		`SELECT id, title, content, author, created_at FROM posts WHERE id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

func (r *PostgresRepository) Create(title, content, author string) (post.Post, error) {
	var p post.Post
	err := r.db.QueryRow(
		`INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id, title, content, author, created_at`,
		title, content, author,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

func (r *PostgresRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
