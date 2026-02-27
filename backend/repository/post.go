// repository層: 「データベース操作の専門家」です。
// SQL（データベースを操作するための言語）を使って、データの「保存・取得・更新・削除」を行います。
// ここにはビジネスルール（NGワードなど）は書かず、純粋なDB操作のみを書きます。
package repository

import (
	"database/sql"

	"github.com/arakou0812/backend/model"
)

// PostRepository: データベース操作を行う本体です。
// *sql.DB（データベースへの接続情報）を保持します。
type PostRepository struct {
	db *sql.DB
}

// NewPostRepository: DB接続情報を受け取って、新しいRepositoryを作ります。
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// GetAll: データベースから「全ての投稿」を取得します。
func (r *PostRepository) GetAll() ([]model.Post, error) {
	// 【1】SQLクエリを実行します（新着順：DESC）。
	rows, err := r.db.Query(`SELECT id, title, content, author, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	// 関数の最後に必ず結果セット(rows)を閉じます。
	defer rows.Close()

	// 【2】結果を1行ずつ読み込んで、スライス(配列)に追加していきます。
	var posts []model.Post
	for rows.Next() {
		var p model.Post
		// DBの各カラムの値を、構造体の各フィールドにマッピングします。
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	// 1件もなかった場合に空のリストを返すための処理です。
	if posts == nil {
		posts = []model.Post{}
	}
	return posts, nil
}

// GetByID: IDを指定して「特定の1件」をデータベースから取得します。
func (r *PostRepository) GetByID(id int) (model.Post, error) {
	var p model.Post
	// $1 はプレースホルダーと呼ばれ、安全に値を埋め込む仕組みです。
	err := r.db.QueryRow(
		`SELECT id, title, content, author, created_at FROM posts WHERE id = $1`, id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Create: データベースに「新しい投稿」を保存（INSERT）します。
func (r *PostRepository) Create(title, content, author string) (model.Post, error) {
	var p model.Post
	// RETURNING を使うことで、保存された直後のIDや作成日時を同時に受け取れます。
	err := r.db.QueryRow(
		`INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id, title, content, author, created_at`,
		title, content, author,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
	return p, err
}

// Delete: データベースから「特定の1件」を削除します。
func (r *PostRepository) Delete(id int) (bool, error) {
	// Exec は、値を返さない操作（DELETE, UPDATE等）に使います。
	result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	// 実際に何行削除されたか（RowsAffected）を確認します。
	rowsAffected, _ := result.RowsAffected()
	// 1行以上削除されていれば成功(true)を返します。
	return rowsAffected > 0, nil
}
