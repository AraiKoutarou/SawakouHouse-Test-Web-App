// persistence層: 「データの永続化（保存）」を担当します。
// 実際にデータベース(PostgreSQL)に対してSQLを発行し、データの保存、取得、更新、削除を行います。
package persistence

import (
	"database/sql"
	"log"

	"github.com/arakou0812/backend/modules/post/domain"
)

// PostgresPostRepository: PostgreSQLデータベースを操作するための部品です。
type PostgresPostRepository struct {
	db *sql.DB // データベースへの接続情報を持ちます。
}

// NewPostRepository: DB接続情報を受け取り、リポジトリを作成します。
func NewPostRepository(db *sql.DB) *PostgresPostRepository {
	return &PostgresPostRepository{db: db}
}

// Migrate: アプリケーションが動くために「最低限必要なテーブル」を自動的に作成します。
// これにより、誰が環境を構築しても同じテーブル構成になります。
func (r *PostgresPostRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id         SERIAL PRIMARY KEY,           -- 自動で増えるID（背番号）
		title      VARCHAR(255) NOT NULL,        -- タイトル（最大255文字、必須）
		content    TEXT NOT NULL,               -- 本文（文字数無制限、必須）
		author     VARCHAR(100) NOT NULL,        -- 投稿者名（最大100文字、必須）
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() -- 投稿日時（自動で現在時刻が入る）
	);`

	// DBにクエリを投げ、テーブルがなければ作ります。
	if _, err := r.db.Exec(query); err != nil {
		log.Fatalf("postsテーブル作成失敗: %v", err)
	}
	log.Println("postsテーブル確認完了")
}

// GetAll: データベースから全ての投稿を取得し、新着順（投稿日時の降順）に並べて返します。
func (r *PostgresPostRepository) GetAll() ([]domain.Post, error) {
	// SELECT文で必要なカラムをすべて指定して取得します。
	rows, err := r.db.Query(`SELECT id, title, content, author, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 使い終わったら確実にリソースを解放します。

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		// DBから取ってきた1行ずつのデータを、プログラムで扱える「Post構造体」に詰め替えます。
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	// 投稿が1つもなかった場合は空のリストを返します。
	if posts == nil {
		posts = []domain.Post{}
	}
	return posts, nil
}

// GetByID: 指定されたIDの投稿を1件だけ取得します。
func (r *PostgresPostRepository) GetByID(id int) (domain.Post, error) {
	var p domain.Post
	// $1 という記法は「プレースホルダー」と呼ばれ、SQLインジェクションという攻撃を防ぐための重要な技術です。
	err := r.db.QueryRow(
		`SELECT id, title, content, author, created_at FROM posts WHERE id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Create: 渡されたタイトル、本文、投稿者名をDBに保存し、自動生成されたIDや日時を含めて返します。
func (r *PostgresPostRepository) Create(title, content, author string) (domain.Post, error) {
	var p domain.Post
	// RETURNING 句を使うことで、INSERT（保存）と同時に「生成されたID」や「作成日時」を取得できます。
	err := r.db.QueryRow(
		`INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id, title, content, author, created_at`,
		title, content, author,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Delete: 指定されたIDの投稿を削除します。
func (r *PostgresPostRepository) Delete(id int) (bool, error) {
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	// 実際に何行のデータが削除されたかを確認します。
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0, nil // 1行以上消えていれば「成功（true）」を返します。
}
