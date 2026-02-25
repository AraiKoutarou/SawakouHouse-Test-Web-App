// repository層: DBへのSQL操作だけを担当する
// ビジネスルールは書かない。「どうやって保存・取得するか」だけを知っている。
package repository

import (
	"database/sql"

	"github.com/arakou0812/backend/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// GetAll: 全投稿を新着順で返す
func (r *PostRepository) GetAll() ([]model.Post, error) {
	rows, err := r.db.Query(`SELECT id, title, content, author, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if posts == nil {
		posts = []model.Post{}
	}
	return posts, nil
}

// GetByID: IDで1件取得。見つからなければ sql.ErrNoRows を返す
func (r *PostRepository) GetByID(id int) (model.Post, error) {
	var p model.Post
	err := r.db.QueryRow(
		`SELECT id, title, content, author, created_at FROM posts WHERE id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Create: 新規投稿を保存して返す
func (r *PostRepository) Create(title, content, author string) (model.Post, error) {
	var p model.Post
	err := r.db.QueryRow(
		`INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id, title, content, author, created_at`,
		title, content, author,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Delete: IDで1件削除。削除対象がなければ false を返す
func (r *PostRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil
}
