// domain層: postモジュールの核心です。
// 他のどのパッケージにも依存しません。
// エンティティ（データの形）とビジネスエラーの定義だけを行います。
package domain

import (
	"errors"
	"time"
)

// Post: 投稿を表すエンティティです。
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// ビジネスエラーの定義。
// service層ではなくdomain層に置くことで、どの層からも参照できます。
var (
	ErrNotFound = errors.New("投稿が見つかりません")
	ErrNGWord   = errors.New("NGワードが含まれています")
)
