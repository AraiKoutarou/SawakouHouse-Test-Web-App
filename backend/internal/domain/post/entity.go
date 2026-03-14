// internal/domain/post/entity.go: ドメインモデル。他パッケージに一切依存しません。
package post

import (
	"errors"
	"time"
)

// Post: 投稿エンティティ。
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// ドメイン例外。
var (
	ErrNotFound = errors.New("投稿が見つかりません")
	ErrNGWord   = errors.New("NGワードが含まれています")
)
